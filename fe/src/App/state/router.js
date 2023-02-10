"use strict"

import Router from '../../lib/router.js';
import routes from './routes.js';
import { curPath } from './store.js';

let router = new Router();
for (let [k, v] of Object.entries(routes)) {
  router.Add(v.path, v.name, v.handler);
}

function moveTo(path) {
  curPath.set(path);

  window.history.pushState(
    {path: path},
    '',
    window.location.origin + path);
}

function moveBack() {
  window.history.back();
}

export default router;
export { moveTo, moveBack };
