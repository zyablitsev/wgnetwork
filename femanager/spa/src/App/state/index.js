import routes from './routes.js';
import router from './router.js';
import { moveTo, moveBack } from './router.js';
import { authToken, curPath, curRoute } from './store.js';

export {
  routes, router,
  moveTo, moveBack,
  authToken, curPath, curRoute,
};
