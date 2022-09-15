function getPageLocale() {
  return $('.header__logo > a').first().attr('href') === '/ru/' ? 'ru' : 'en';
}

function getHue4Version(version) {
  let result = 0
  for (let i = 0; i < version.length; ++i)
    result = Math.imul(31, result) + version.charCodeAt(i)
  result = Math.abs(result)
  return result
}

function getEditionName(edition) {
  switch (edition) {
    case 'fe':
      return 'FE (Flant Edition)';
    case 'ee':
      return 'EE (Enterprise Edition)';
    case 'ce':
      return 'CE (Community Edition)';
    default:
      return '';
  }
}

function getTitle() {
  return getPageLocale() === 'ru' ? 'Дата появления версии на канале обновлений' : 'Date the version appeared on the release channel';
}

function formatDate(date) {
  return new Intl.DateTimeFormat(getPageLocale() === 'ru' ? 'ru-RU' : 'en-US', {
      weekday: 'long',
      day: 'numeric',
      month: 'long'
  }).format(date);
}

async function showReleaseChannelStatus(apiURL) {
  const channels = ['a', 'b', 'ea', 's', 'rs'];
  const channelCodes = {
    "Alpha": 'a',
    "Beta": 'b',
    "Early Access": 'ea',
    "Stable": 's',
    "Rock Solid": 'rs' };
  const editions = ['fe', 'ee', 'ce'];

  const root = document.querySelector('.releases-page__table--content');
  const table = document.createElement('table');
  const thead = document.createElement('thead');
  const trHead = document.createElement('tr');
  const tbody = document.createElement('tbody');

  await fetch(apiURL, {
      headers: {
        'Accept': 'application/json'
      },
    })
    .then(respose => respose.json())
    .then(data => {
      for (const channelItem in channelCodes) {
        const trBody = document.createElement('tr');
        const channel = document.createElement('td');
        channel.classList.add('channel-name');
        channel.innerText = channelItem[0].toUpperCase() + channelItem.slice(1);
        trBody.append(channel)
        for (const edition of editions) {
           if (edition === 'fe') continue;
           const itemData = data.releases[channelCodes[channelItem] + '-' + edition]
           const rawItem = document.createElement('td');
           const linkBlock = document.createElement('p');
           linkBlock.classList.add('link');
           const link = document.createElement('a');
           linkBlock.append(link);
           const date = new Date(Date.parse(itemData['date']))
           link.href = `../${itemData['version']}/`;
           link.innerText = itemData['version'].replace(/^v/,'');
           // rawItem.innerText = `(${formatDate(date)})`;
           const dateItem = document.createElement('p');
           dateItem.classList.add('date');
           dateItem.innerText = `(${formatDate(date)})`;
           dateItem.setAttribute('title', getTitle());
           rawItem.append(dateItem);
           // rawItem.prepend(document.createElement('br'));
           rawItem.prepend(linkBlock);
           rawItem.style = "background-color: hsl(" + getHue4Version(itemData['version']) + ", 50%, 85%);";
           trBody.append(rawItem)
        }
        tbody.append(trBody);
      }
    })

  let th = document.createElement('th');

  th.innerText = getPageLocale() === 'ru' ? 'Канал обновлений' : 'Release channel';
  trHead.append(th);
  thead.append(trHead);

  for (const edition of editions) {
    if (edition === 'fe') continue;
    th = document.createElement('th');
    th.innerText = getEditionName(edition);
    trHead.append(th);
  }

  table.append(thead);
  table.append(tbody);
  root.append(table);
}
async function showReleaseChannelStatusHorizontal(apiURL) {
  const channels = ['a', 'b', 'ea', 's', 'rs'];
  const channelCodes = {
    "Alpha": 'a',
    "Beta": 'b',
    "Early Access": 'ea',
    "Stable": 's',
    "Rock Solid": 'rs' };
  const editions = ['fe', 'ee', 'ce'];

  const root = document.querySelector('.releases-page__table--content');
  const table = document.createElement('table');
  const thead = document.createElement('thead');
  const trHead = document.createElement('tr');
  const tbody = document.createElement('tbody');

  await fetch(apiURL, {
      headers: {
        'Accept': 'application/json'
      },
    })
    .then(respose => respose.json())
    .then(data => {
      for (const edition of editions) {
        if (edition === 'fe') continue;
        const trBody = document.createElement('tr');
        const editionItem = document.createElement('td');
        editionItem.classList.add('edition-name');
        editionItem.innerText = getEditionName(edition);
        trBody.append(editionItem)
        for (const channelItem in channelCodes) {
           const itemData = data.releases[channelCodes[channelItem] + '-' + edition]
           const rawItem = document.createElement('td');
           const linkBlock = document.createElement('p');
           linkBlock.classList.add('link');
           const link = document.createElement('a');
           linkBlock.append(link);
           const date = new Date(Date.parse(itemData['date']))
           link.href = `../${itemData['version']}/`;
           link.innerText = itemData['version'].replace(/^v/,'');
           // rawItem.innerText = `(${formatDate(date)})`;
           const dateItem = document.createElement('p');
           dateItem.classList.add('date');
           dateItem.innerText = `(${formatDate(date)})`;
           dateItem.setAttribute('title', getTitle());
           rawItem.append(dateItem);
           // rawItem.prepend(document.createElement('br'));
           rawItem.prepend(linkBlock);
           rawItem.style = "background-color: hsl(" + getHue4Version(itemData['version']) + ", 50%, 85%);";
           trBody.append(rawItem)
        }
        tbody.append(trBody);
      }
    })

  let th = document.createElement('th');

  th.innerText = getPageLocale() === 'ru' ? 'Канал обновлений' : 'Release channel';
  trHead.append(th);
  thead.append(trHead);

  for (const channelItem in channelCodes) {
    th = document.createElement('th');
    th.innerText = channelItem[0].toUpperCase() + channelItem.slice(1);
    trHead.append(th);
  }

  table.append(thead);
  table.append(tbody);
  root.append(table);
}

document.addEventListener("DOMContentLoaded", function () {
  const apiURL = 'https://flow.deckhouse.io/releases';
  showReleaseChannelStatusHorizontal(apiURL)
    .then(() =>  {
      $('.releases-page__table--content').addClass('active');
      })
    .catch((reason) => {
      $('.releases-page__loadblock.failed').addClass('active');
      console.log(`Failed to fetch release channel data from ${apiURL}.`, reason)
      })
    .finally( () => {
      $('.releases-page__loadblock.progress').removeClass('active')
      });

  // .catch((reason) => {
  //   $('.releases-page__onload.failed').addClass('active');
  //   console.log(`Failed to call showReleaseChannelStatus.`, reason)}
  // )
});
