/** @format */

import { EventBus } from './eventbus';

export default {
  log(v) {
    console.log(v);
  },

  catchRestError(err) {
    if (err.response) {
      var response = err.response;
      var request = response.request;
      EventBus.$emit(
        'main-error',
        `Request failed with status code ${
          response.status
        } and message: ${JSON.stringify(response.data)}`
      );
    } else if (err.request) {
      EventBus.$emit(
        'main-error',
        `Request failed with no response from the back end.`
      );
    } else {
      EventBus.$emit('main-error', `Request failed: ${err.message}`);
    }
  },

  copySLToClipboard(subUrl) {
    return new Promise((resolve, reject) => {
      var id = 'hidden-clipboard-area';
      var existsTextarea = document.getElementById(id);
      if (!existsTextarea) {
        var textarea = document.createElement('textarea');
        textarea.id = id;
        document.querySelector('body').appendChild(textarea);
        existsTextarea = document.getElementById(id);
      }
      existsTextarea.value = window.location.origin + '/' + subUrl;
      existsTextarea.select();
      var status = document.execCommand('copy');
      if (!status) {
        reject('Could not copy shortlink to clipboard.');
      } else {
        resolve();
      }
    });
  },
};
