<!-- ----------- -->
<!--  TEMPLATES  -->
<!-- ----------- -->
<template>
  <div id="app">
    <!-- ERROR ALERT -->
    <b-alert :show="alert.visible" :variant="alert.type">
      {{ alert.msg }}
    </b-alert>

    <!-- ENTRY LIST -->
    <Entry v-for="sl in shortlinks" 
      :key="sl.id"
      :id="sl.id"
      :rootlink="sl.root_link"
      :shortlink="sl.short_link"
      :created="sl.created"
      :accesses="sl.accesses"
      :edited="sl.edited"
    />

    <!-- ADD BUTTON -->
    <a v-b-modal.modal-add v-if="addButtonVisible" class="add text-white">+</a>

    <!-- LOGIN MODAL -->
    <b-modal 
      id="modal-login" 
      class="text-dark" 
      title="Login"
      @hide="login(loginModal.tbToken)"
    >
      <b-alert :show="loginModal.showWrongCredentials" variant="danger">
        Wrong login credentials.
      </b-alert>
      <b-form-group
        id="groupLoginToken"
        label="Login token:"
        label-for="tbLoginToken"
      >
        <b-form-input
          id="tbShortLink"
          ref="tbShortLink"
          type="password"
          placeholder=""
          v-model="loginModal.tbToken"
          @keyup.enter.native="$bvModal.hide('modal-login')"
        ></b-form-input>
      </b-form-group>
    </b-modal>

    <!-- ADD MODAL -->
    <b-modal 
      id="modal-add" 
      class="text-dark" 
      title="Add Short Link"
      @ok="createShortLink(addModal.tbRootLink, addModal.tbShortLink)"  
    >
      <b-form-group
        id="groupShortLink"
        label="Short link"
        label-for="tbShortLink"
        description="Leave empty for random short link identifier."
      >
        <b-form-input
          id="tbShortLink"
          type="text"
          placeholder=""
          v-model="addModal.tbShortLink"
        ></b-form-input>
      </b-form-group>
      <b-form-group
        id="groupRootLink"
        label="Root link"
        label-for="tbRootLink"
        description="URI the link will redirect to."
      >
        <b-form-input
          id="tbRootLink"
          type="text"
          required
          placeholder=""
          v-model="addModal.tbRootLink"
        ></b-form-input>
      </b-form-group>
    </b-modal>

    <!-- EDIT MODAL -->
    <b-modal 
      id="modal-edit" 
      class="text-dark" 
      title="Edit Short Link"
      @ok="modifyShortLink(editModal.id, editModal.tbRootLink, editModal.tbShortLink)"  
    >
      <b-form-group
        id="groupShortLink"
        label="Short link"
        label-for="tbShortLink"
        description="The short link identifier."
      >
        <b-form-input
          id="tbShortLink"
          type="text"
          placeholder=""
          v-model="editModal.tbShortLink"
        ></b-form-input>
      </b-form-group>
      <b-form-group
        id="groupRootLink"
        label="Root link"
        label-for="tbRootLink"
        description="URI the link will redirect to."
      >
        <b-form-input
          id="tbRootLink"
          type="text"
          required
          placeholder=""
          v-model="editModal.tbRootLink"
        ></b-form-input>
      </b-form-group>
    </b-modal>

    <!-- DELETE MODAL -->
    <b-modal 
      id="modal-delete" 
      hide-footer 
      class="text-dark" 
      title="Delete Short Link"
    >
      <p>Do you really want to delete the following short link?</p>
      <table class="mod-delete">
        <tr>
          <th>ID</th>
          <td>{{ deleteModal.id }}</td>
        </tr>
        <tr>
          <th>Short Link</th>
          <td>{{ deleteModal.shortLink }}</td>
        </tr>
        <tr>
          <th>Root URI</th>
          <td>{{ deleteModal.rootLink }}</td>
        </tr>
      </table>
      <div class="text-right">
        <b-button variant="secondary" class="mt-3 mr-2" @click="$bvModal.hide('modal-delete')">Cancel</b-button>
        <b-button variant="danger" class="mt-3" @click="deleteShortLink(deleteModal.id)">Delete</b-button>
      </div>
    </b-modal>
  </div>
</template>

<!-- -------- -->
<!--  SCRIPT  -->
<!-- -------- -->
<script>
import Entry from './components/Entry.vue'
import rest from './js/rest';
import utils from './js/utils';
import { EventBus } from './js/eventbus';

export default {
  name: 'app',
  components: {
    Entry,
  },

  data() {
    return {
      shortlinks: [],

      alert: {
        msg: '',
        type: 'danger',
        visible: false,
      },

      addModal: {
        tbShortLink: null,
        tbRootLink: null,
      },

      editModal: {
        tbShortLink: null,
        tbRootLink: null,
        id: null,
      },

      deleteModal: {
        shortLink: null,
        rootLink: null,
        id: null,
      },

      loginModal: {
        tbToken: null,
        showWrongCredentials: false,
      },

      addButtonVisible: false,
    };
  },

  methods: {
    refetchData() {
      rest.getShortlinks().then(res => {
        this.addButtonVisible = true;
        this.shortlinks = res.data.results;
      }).catch((err) => {
        if (err.response && err.response.status == 401) {
          this.$bvModal.show('modal-login');
        } else {
          utils.catchRestError(err);
        }
      });
    },

    login(token) {
      var showModal = (err) => {
        this.loginModal.tbToken = '';
        if (err.response && err.response.status == 401) {
          this.loginModal.showWrongCredentials = true;
        }
        setTimeout(() => 
          this.$bvModal.show('modal-login'), 10);
      };

      if (!token || token.length < 1) {
        showModal({});
      } else {
        rest.login(token)
          .catch((err) => showModal(err))
          .then((res) => this.refetchData());
      }
    },

    test(e) {
      console.log(e);
    },

    createShortLink(root, short) {
      rest.createShortLink(root, short).then((res) => {
        this.refetchData();
        EventBus.$emit('main-info', 
          `Successfully created short link ${res.data.short_link} (ID: ${res.data.id})`);        
      }).catch(utils.catchRestError);
    },

    modifyShortLink(id, root, short) {
      rest.modifyShortLink(id, root, short).then((res) => {
        this.refetchData();
        EventBus.$emit('main-info', 
          `Successfully updated short link ${res.data.short_link} (ID: ${res.data.id})`);
      }).catch(utils.catchRestError);
    },

    deleteShortLink(id) {
      rest.deleteShortLink(id).then((res) => {
        this.refetchData();
        EventBus.$emit('main-info', 
          `Successfully deleted short link (ID was ${res.data.id}).`);
        this.$bvModal.hide('modal-delete')
      }).catch(utils.catchRestError);
    }
  },

  created: function() {
    EventBus.$on('main-error', (msg) => {
      this.alert.visible = true;
      this.alert.type = 'danger';
      this.alert.msg = msg;
      setTimeout(() => {
        this.alert.visible = false;
      }, 5000);
    });

    EventBus.$on('main-info', (msg) => {
      this.alert.visible = true;
      this.alert.type = 'success';
      this.alert.msg = msg;
      setTimeout(() => {
        this.alert.visible = false;
      }, 5000);
    });

    EventBus.$on('main-edit', (data) => {
      this.editModal.tbShortLink = data.shortlink;
      this.editModal.tbRootLink = data.rootlink;
      this.editModal.id = data.id;
      this.$bvModal.show('modal-edit');
    });

    EventBus.$on('main-delete', (data) => {
      this.deleteModal.shortLink = data.shortlink;
      this.deleteModal.rootLink = data.rootlink;
      this.deleteModal.id = data.id;
      this.$bvModal.show('modal-delete');
    });

    this.refetchData();
  }, 
}
</script>

<!-- ------- -->
<!--  STYLE  -->
<!-- ------- -->
<style>

#app {
  width: 100%;
}

body {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  background-color: #263238;
  padding: 20px;
}

a.add {
  position: fixed;
  bottom: 30px;
  right: 30px;
  border-radius: 50%;
  color: white;
  background-color: #1565C0;
  font-size: 45px;
  padding: 0px 21px;
  margin: 0px;
  box-shadow: 0px 0px 75px 23px rgba(0,0,0,0.75);
  cursor: pointer;
  transition: all .25s ease;
}

a.add:hover {
  transform: scale(1.1);
}

b-modal {
  color: black !important;
}

table.mod-delete th {
  padding-right: 50px;
}

#hidden-clipboard-area {
    position: fixed;
    top: 0;
    left: 0;
    width: 1px;
    height: 1px;
    padding: 0;
    border: none;
    outline: none;
    box-shadow: none;
    background: transparent;
}

</style>
