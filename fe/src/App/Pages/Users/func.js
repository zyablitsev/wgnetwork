"use strict"

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
      role: users[i]['is_manager'] ? 'manager' : 'member',
      devices: users[i]['devices'].length
    };
  }

  let result = {
    session: check.session,
    users: users,
  }

  return Promise.resolve(result);
}

export {
  getUsers,
};
