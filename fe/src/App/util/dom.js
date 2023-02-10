"use strict"

function setHtmlClass(...cssClass) {
  let root = document.getElementsByTagName('html');
  if (root !== "undefined" && root.length > 0) {
    root = root[0];
  }
  root.className = ''; // remove all classes
  root.classList.add(...cssClass);
}

function setBodyClass(...cssClass) {
  document.body.className = ''; // remove all classes
  document.body.classList.add(...cssClass);
}

export {
  setHtmlClass,
  setBodyClass
};
