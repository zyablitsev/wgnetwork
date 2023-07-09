"use strict"

function unsecuredCopyToClipboard(text) {
  const textArea = document.createElement("textarea");
  textArea.value = text;
  document.body.appendChild(textArea);
  textArea.style ="white-space: pre";
  textArea.focus();
  textArea.select();
  try {
    document.execCommand('copy');
  } catch (err) {
    console.error('Unable to copy to clipboard', err);
  }
  document.body.removeChild(textArea);
}

function copyToClipboard(event) {
  event.preventDefault();
  let elText = this.parentElement.parentElement.lastChild;
  let content = elText.textContent.trim();
  if (navigator.clipboard) {
    navigator.clipboard.writeText(content);
  } else {
    unsecuredCopyToClipboard(content);
  }
  let elNotice = this.parentElement.firstChild;

  // TODO: rework and fix
  elNotice.classList.remove('invisible');
  elNotice.classList.remove('opacity-0');
  elNotice.classList.add('opacity-100');
  setTimeout(() => {
    elNotice.classList.remove('opacity-100');
    elNotice.classList.add('opacity-0');
    elNotice.classList.add('invisible');
  }, 2000);
};

export {
  copyToClipboard,
};
