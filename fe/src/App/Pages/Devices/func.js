"use strict"

async function getDevices(client, session) {
  let check = {};
  try {
    check = await client.Fetch('auth/session/check', {}, session);
  } catch (err) {
    return Promise.reject(err);
  }

  let users = [];
  let devices = [];
  try {
    [users, devices] = await Promise.all([
      client.Fetch('manager/users', {}, session),
      client.Fetch('manager/devices', {}, session)]);
  } catch (err) {
    return Promise.reject(err);
  }

  let usersMap = {};
  for (let i = 0; i < users.length; i++) {
    usersMap[users[i]['uuid']] = users[i]['name'];
  }

  for (let i = 0; i < devices.length; i++) {
    let username = usersMap[devices[i]['user_uuid']];

    let ip = devices[i]['ipnetwork'].split('/');
    ip = (ip.length > 0) ? ip[0] : '';

    devices[i] = {
      ipnetwork: devices[i]['ipnetwork'],
      label: devices[i]['label'],
      username: username,
      ip: ip
    };
  }

  let result = {
    session: check.session,
    users: users,
    devices: devices
  }

  return Promise.resolve(result);
}

export {
  getDevices,
};
