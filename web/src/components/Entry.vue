<template>
  <b-container fluid class="entry mb-3 p-2 rounded text-white">
    <b-row>
      <b-col>
        <h5 class="mb-0">{{ shortlink }}</h5>
        <div>
          <a :href="rootlink" target="_blank" class="root mt-1">{{ shorten(rootlink, 80) }}</a>
        </div>
        <b-badge variant="secondary">#{{ id }}</b-badge>
        <b-badge variant="primary ml-2">{{ accesses }} Clicks</b-badge>
      </b-col>
      <b-col cols="2" class="text-right">
        <b-button 
          variant="warning" 
          class="mr-2"
          @click="copy"
        >🔗</b-button>
        <b-button 
          variant="info" 
          class="mr-2"
          @click="edit"
        >✏️</b-button>
        <b-button 
          variant="danger"
          @click="del"
        >🗑️</b-button>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
/** @format */

import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';
import { EventBus } from '../js/eventbus';
import utils from '../js/utils';

export default {
  name: 'Entry',
  
  props: {
    id: Number,
    rootlink: String,
    shortlink: String,
    created: String,
    accesses: Number,
    edited: String,
  },

  methods: {
    shorten(txt, maxlen) {
      maxlen = maxlen - 3;
      if (txt.length > maxlen) {
        txt = txt.substr(0, maxlen) + '...';
      }
      return txt;
    },

    edit() {
      EventBus.$emit('main-edit', {
        id: this.id,
        shortlink: this.shortlink,
        rootlink: this.rootlink,
      });
    },

    del() {
      EventBus.$emit('main-delete', {
        id: this.id,
        shortlink: this.shortlink,
        rootlink: this.rootlink,
      });
    },

    copy() {
      utils.copySLToClipboard(this.shortlink)
        .then(() => EventBus.$emit('main-info', 'Short link copied to clip board.'))
        .catch((err) => EventBus.$emit('main-error', err));
    },
  },
};
</script>

<style scoped>
.entry {
  background-color: #37474F;
}

a.root {
  color: #BDBDBD;
  margin: 0px;
  font-size: 14px;
}
</style>
