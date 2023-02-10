"use strict"

async function getData(client, s) {
  let check = {};
  try {
    check = await client.Fetch('auth/session/check', {}, s);
  } catch (err) {
    return Promise.reject(err);
  }

  let wgCfg = {};
  let ipSet = [];
  try {
    [wgCfg, ipSet] = await Promise.all([
      client.Fetch('manager/wg/cfg', {}, s),
      client.Fetch('manager/trust/ipset', {}, s)]);
  } catch (err) {
    return Promise.reject(err);
  }

  let result = {
    session: check['session'],
    wgCfg: wgCfg,
    ipSet: ipSet
  };

  return Promise.resolve(result);
}

async function ipsetAdd(client, s, ip) {
  let r = {};
  try {
    let p = {'ip': ip};
    r = await client.Fetch('manager/trust/ipset/add', p, s);
  } catch (err) {
    return Promise.reject(err);
  }

  let ipset = [];
  try {
    ipset = await client.Fetch('manager/trust/ipset', {}, s);
  } catch (err) {
    return Promise.reject(err);
  }

  return Promise.resolve(ipset);
}

async function ipsetRemove(client, s, ip) {
  let r = {};
  try {
    let p = {'ip': ip};
    r = await client.Fetch('manager/trust/ipset/remove', p, s);
  } catch (err) {
    return Promise.reject(err);
  }

  let ipset = [];
  try {
    ipset = await client.Fetch('manager/trust/ipset', {}, s);
  } catch (err) {
    return Promise.reject(err);
  }

  return Promise.resolve(ipset);
}

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

function handleClipboard(event) {
  event.preventDefault();
  let elText = this.parentElement.parentElement.lastChild;
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

export {
  getData,
  ipsetAdd,
  ipsetRemove,
  handleClipboard
};
