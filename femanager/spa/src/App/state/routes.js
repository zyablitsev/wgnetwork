"use strict"

import Server from '../Pages/Server/Server.svelte';
import Users from '../Pages/Users/Users.svelte';
import User from '../Pages/User/User.svelte';
import UserCreate from '../Pages/User/UserCreate.svelte';
import UserEdit from '../Pages/User/UserEdit.svelte';
import Devices from '../Pages/Devices/Devices.svelte';
import Device from '../Pages/Device/Device.svelte';
import DeviceCreate from '../Pages/Device/DeviceCreate.svelte';
import DeviceEdit from '../Pages/Device/DeviceEdit.svelte';
import SignIn from '../Pages/SignIn/SignIn.svelte';

let routes = {
  "server": {"path": "/", "name": "server", "handler": Server},
  "users": {"path": "/users", "name": "users", "handler": Users},
  "user": {"path": "/user/:uuid", "name": "user", "handler": User},
  "user_create": {"path": "/user/create", "name": "user_create", "handler": UserCreate},
  "user_edit": {"path": "/user/:uuid/edit", "name": "user_edit", "handler": UserEdit},
  "devices": {"path": "/devices", "name": "devices", "handler": Devices},
  "device": {"path": "/device/:ip", "name": "device", "handler": Device},
  "device_create": {"path": "/device/create", "name": "device_create", "handler": DeviceCreate},
  "device_edit": {"path": "/device/:ip/edit", "name": "device_edit", "handler": DeviceEdit},
  "signIn": {"path": "/signin", "name": "signin", "handler": SignIn},
}

export default routes;
