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

async function getUser(client, session, uuid) {
  let check = {};
  try {
    check = await client.Fetch('auth/session/check', {}, session);
  } catch (err) {
    return Promise.reject(err);
  }

  let user = {};
  let params = {'uuid': uuid};
  try {
    user = await client.Fetch('manager/user', params, session);
  } catch (err) {
    return Promise.reject(err);
  }

  user = {
    uuid: user['uuid'],
    name: user['name'],
    isManager: user['is_manager'],
    role: user['is_manager'] ? 'manager' : 'member',
  };

  let result = {
    session: check.session,
    user: user
  }

  return Promise.resolve(result);
}

async function getUserWithDevices(client, session, uuid) {
  let check = {};
  try {
    check = await client.Fetch('auth/session/check', {}, session);
  } catch (err) {
    return Promise.reject(err);
  }

  let user = {};
  let devices = [];
  try {
    let params = {'uuid': uuid};
    [user, devices] = await Promise.all([
      client.Fetch('manager/user', params, session),
      client.Fetch('manager/devices', {}, session)]);
  } catch (err) {
    return Promise.reject(err);
  }

  let devicesMap = {};
  for (let i = 0; i < devices.length; i++) {
    let ip = devices[i]['ipnetwork'].split("/");
    ip = (ip.length > 0) ? ip[0] : '';

    if (ip.length < 1) {
      continue;
    }

    devicesMap[ip] = devices[i];
  }

  devices = [];
  for (let i = 0; i < user['devices'].length; i++) {
    devices.push({
      ipnetwork: devicesMap[user['devices'][i]]['ipnetwork'],
      ip: user['devices'][i],
      label: devicesMap[user['devices'][i]]['label']
    });
  }

  user = {
    uuid: user['uuid'],
    name: user['name'],
    isManager: user['is_manager'],
    role: user['is_manager'] ? 'manager' : 'member',
    devices: devices,
    key: user['key'],
    provisionUri: user['provision_uri']
  };

  let result = {
    session: check.session,
    user: user,
    devices: devices
  }

  return Promise.resolve(result);
}

function formValidate(username) {
  if (username.length < 1 || username.length > 64) {
    return false;
  }

  return true;
}

export {
  checkSession,
  getUser,
  getUserWithDevices,
  formValidate,
};
