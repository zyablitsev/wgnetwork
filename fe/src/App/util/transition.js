"use strict"

import { linear } from 'svelte/easing';

function transitionCfg(node, params) {
  const isOut = params.isOut ? true : false;
  const isIn = !isOut;

  const base = params.base ? params.base.split(' ') : [];
  const from = params.from ? params.from.split(' ') : [];
  const to = params.to ? params.to.split(' ') : [];

  if (base.length > 0 || from.length > 0) {
    node.classList.add(...[...base, ...from]);
  }

  let started = false;

  return {
    delay: params.delay || 0,
    duration: params.duration || 200,
    easing: params.easing || linear,
    tick: (t, u) => {
      if ((t === 1 && isIn) || (t === 0 && isOut)) {
        if (base.length > 0) {
          node.classList.remove(...base);
        }

        return;
      }

      if (started) {
        return;
      }

      if ((t > 0 && isIn) || (t < 1 && isOut)) {
        started = true;

        if (from.length > 0) {
          node.classList.remove(...from);
        }

        if (to.length > 0) {
          node.classList.add(...to);
        }
      }
    },
  };
}

export { transitionCfg };
