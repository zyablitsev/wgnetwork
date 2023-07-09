import { writable, derived } from 'svelte/store';

export const authToken = writable('');
export const curPath = writable('');
export const curRoute = derived(
  [curPath, authToken],
  ([$curPath, $authToken]) => {
    return {"token": $authToken, "path": $curPath};
  }
);
