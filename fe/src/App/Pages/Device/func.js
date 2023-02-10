"use strict"

async function checkSession(client, session) {
  let check = {};
  try {
    check = await client.Fetch('auth/session/check', {}, session);
  } catch (err) {
    return Promise.reject(err);
  }

  let result = {
    session: check.session
  }

  return Promise.resolve(result);
}

async function getDevice(client, session, ip) {
  let check = {};
  try {
    check = await client.Fetch('auth/session/check', {}, session);
  } catch (err) {
    return Promise.reject(err);
  }

  let device = {};
  let params = {'ip': ip};
  try {
    device = await client.Fetch('manager/device', params, session);
  } catch (err) {
    return Promise.reject(err);
  }

  let wgcfg = {
    wgDeviceInet: device['wg_device_inet'],
    wgDevicePort: device['wg_device_port'],
    wgDevicePubKey: device['wg_device_pubkey'],
    wgDeviceAllowedIPs: device['wg_device_allowed_ips'],
    wgDeviceDNS: ['8.8.8.8', '8.8.4.4'],

    wgServerInet: device['wg_server_inet'],
    wgServerIPNet: device['wg_server_ipnet'],
    wgServerPort: device['wg_server_port'],
    wgServerPubKey: device['wg_server_pubkey'],

    serverWanIP: device['server_wanip'],
  };

  ip = device['wg_device_inet'].split('/');
  ip = (ip.length > 0) ? ip[0] : '';

  device = {
    ip: ip,
    label: device['label'],
    wanForward: device['wan_forward'],
    user: {
      uuid: device['user_uuid'],
      name: device['user_name']
    }
  }

  let result = {
    session: check.session,
    device: device,
    wgcfg: wgcfg
  }

  return Promise.resolve(result);
}

async function getUsers(client, session) {
  let check = {};
  try {
    check = await client.Fetch('auth/session/check', {}, session);
  } catch (err) {
    return Promise.reject(err);
  }

  let users = [];
  try {
    users = await client.Fetch('manager/users', {}, session);
  } catch (err) {
    return Promise.reject(err);
  }

  for (let i = 0; i < users.length; i++) {
    users[i] = {
      name: users[i]['name'],
      uuid: users[i]['uuid'],
    };
  }

  let result = {
    session: check.session,
    users: users,
  }

  return Promise.resolve(result);
}

function formValidate(label, uuid) {
  if (label.length < 1 || label.length > 64) {
    return false;
  }

  if (uuid.length < 1) {
    return false;
  }

  return true;
}

export {
  checkSession,
  getDevice,
  getUsers,
  formValidate,
};
