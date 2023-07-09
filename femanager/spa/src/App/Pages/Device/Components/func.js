"use strict"

function unsecuredCopyToClipboard(text) {
  const textArea = document.createElement("textarea");
  textArea.value = text;
  document.body.appendChild(textArea);
  textArea.style ="white-space: pre";
  textArea.focus();
  textArea.select();
  try {
    document.execCommand('copy');
  } catch (err) {
    console.error('Unable to copy to clipboard', err);
  }
  document.body.removeChild(textArea);
}

function copyToClipboard(event) {
  event.preventDefault();
  let elText = this.parentElement.parentElement.parentElement.lastChild;
  let content = elText.textContent.trim();
  if (navigator.clipboard) {
    navigator.clipboard.writeText(content);
  } else {
    unsecuredCopyToClipboard(content);
  }
  let elNotice = this.parentElement.firstChild;

  // TODO: rework and fix
  elNotice.classList.remove('invisible');
  elNotice.classList.remove('opacity-0');
  elNotice.classList.add('opacity-100');
  setTimeout(() => {
    elNotice.classList.remove('opacity-100');
    elNotice.classList.add('opacity-0');
    elNotice.classList.add('invisible');
  }, 2000);
};

function download() {
  event.preventDefault();

  let elText = this.parentElement.parentElement.parentElement.lastChild;
  let content = elText.textContent.trim();

  let filename = 'wireguard.txt';
  let element = document.createElement('a');
  element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(content));
  element.setAttribute('download', filename);

  element.style.display = 'none';
  document.body.appendChild(element);

  element.click();

  document.body.removeChild(element);
}

function excludePrivateNetworks(wgipnet, ipnets) {
  if (ipnets.length != 1) {
    return ipnets
  }

  let size = ipnets[0].split('/');
  size = (size.length > 1) ? size[1] : -1;
  if (size < 0) {
    return [wgipnet]
  } else if (size > 0) {
    return ipnets
  }

  ipnets = [
    wgipnet,
    '1.0.0.0/8',
    '2.0.0.0/8',
    '3.0.0.0/8',
    '4.0.0.0/6',
    '8.0.0.0/7',
    '11.0.0.0/8',
    '12.0.0.0/6',
    '16.0.0.0/4',
    '32.0.0.0/3',
    '64.0.0.0/2',
    '128.0.0.0/3',
    '160.0.0.0/5',
    '168.0.0.0/6',
    '172.0.0.0/12',
    '172.32.0.0/11',
    '172.64.0.0/10',
    '172.128.0.0/9',
    '173.0.0.0/8',
    '174.0.0.0/7',
    '176.0.0.0/4',
    '192.0.0.0/9',
    '192.128.0.0/11',
    '192.160.0.0/13',
    '192.169.0.0/16',
    '192.170.0.0/15',
    '192.172.0.0/14',
    '192.176.0.0/12',
    '192.192.0.0/10',
    '193.0.0.0/8',
    '194.0.0.0/7',
    '196.0.0.0/6',
    '200.0.0.0/5',
    '208.0.0.0/4'];

  return ipnets
}

export {
  copyToClipboard,
  download,
  excludePrivateNetworks,
};
