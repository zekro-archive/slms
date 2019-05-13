import axios from "axios";

const ROOTURI = "";
// const ROOTURI = "http://localhost:8080";
const HEADERS = {
  // Authorization: "Basic test"
};

axios.defaults.withCredentials = true;

export default {
  getShortlinks() {
    return axios({
      method: "GET",
      url: ROOTURI + "/api/shortlinks",
      headers: HEADERS
    });
  },

  login(token) {
    return axios({
      method: "POST",
      url: ROOTURI + "/api/login",
      headers: {
        Authorization: "Basic " + token
      }
    });
  },

  createShortLink(root, short) {
    return axios({
      method: "POST",
      url: ROOTURI + "/api/shortlinks",
      headers: HEADERS,
      data: {
        root_link: root,
        short_link: short
      }
    });
  },

  modifyShortLink(id, root, short) {
    return axios({
      method: "POST",
      url: ROOTURI + "/api/shortlinks/" + id,
      headers: HEADERS,
      data: {
        root_link: root,
        short_link: short
      }
    });
  },

  deleteShortLink(id) {
    return axios({
      method: "DELETE",
      url: ROOTURI + "/api/shortlinks/" + id,
      headers: HEADERS
    });
  }
};
