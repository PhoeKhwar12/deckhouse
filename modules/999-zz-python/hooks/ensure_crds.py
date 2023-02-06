#!/usr/bin/env python3
#
# Copyright 2023 Flant JSC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import json
import os

import yaml
from deckhouse import hook
from kubernetes import client
from kubernetes import config as kube_config

# We expect structure with possible subdirectories like this:
# modules/
#   987-your-module-name/
#       crds/
#           crd1.yaml
#           crd2.yaml
#           subdir/
#               crd3.yaml
#       hooks/
#           ensure_crds.py # this file


config = """
configVersion: v1
onStartup: 5
"""


def main():
    hook.run(handle, config=config)


from pprint import pprint


def zzz(s):
    print(f"ZZZ: {s}")


def handle(ctx: hook.Context):
    zzz("Starting handle")
    kube_config.load_incluster_config()
    zzz("Loaded kubeconfig")
    ext_api = client.ApiextensionsV1Api()
    zzz("Created ext_api")

    for crd in iter_manifests(find_crds_root(__file__)):
        try:
            # If Webhook Handler has a conversion webhook for a CRD, it adds '.spec.conversion' to
            # the CRD dynamically. If we blindly re-create the CRD, we will lose the conversion
            # webhook configuration, and conversions will stop working. So we need to read the
            # existing CRD to preserve '.spec.conversion' field.
            # https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/ApiextensionsV1Api.md#read_custom_resource_definition

            name = crd["metadata"]["name"]
            zzz(f"Reading CRD {name}")

            # We just want to put JSON into ctx.kubrentes collector, so we use _preload_content=False to
            # avoid inner library types.
            existing_crd_json = ext_api.read_custom_resource_definition(
                name=name, _preload_content=False
            ).read()
            zzz(f"Read CRD {name}")
            existing_crd = json.loads(existing_crd_json)
            crd["spec"]["conversion"] = existing_crd["spec"]["conversion"]

        except client.rest.ApiException as e:
            if e.status == 404:
                # CRD does not exist, create it
                pass
            else:
                raise e

        zzz("CRD Output")
        pprint(crd)
        ctx.kubernetes.create_or_update(crd)


def iter_manifests(root_path: str):
    if not os.path.exists(root_path):
        return

    for dirpath, dirnames, filenames in os.walk(top=root_path):
        for filename in filenames:
            if not filename.endswith(".yaml"):
                # Wee only seek manifests
                continue
            if filename.startswith("doc-"):
                # Skip dedicated doc yamls, common for Deckhouse internal modules
                continue

            crd_path = os.path.join(dirpath, filename)
            with open(crd_path, "r", encoding="utf-8") as f:
                for manifest in yaml.safe_load_all(f):
                    if manifest is None:
                        continue
                    yield manifest

        for dirname in dirnames:
            subroot = os.path.join(dirpath, dirname)
            for manifest in iter_manifests(subroot):
                yield manifest


def find_crds_root(hookpath):
    hooks_root = os.path.dirname(hookpath)
    module_root = os.path.dirname(hooks_root)
    crds_root = os.path.join(module_root, "crds")
    return crds_root


if __name__ == "__main__":
    main()
