# ⚠ Outdated

This project is not supported anymore. Please visit the successor of this project:  
https://github.com/vctr-sls

---

div align="center">
     <!-- <img src="https://zekro.de/src/go_chat_logo.png" width="400"/> -->
     <h1>~ SLMS ~</h1>
     <strong>Short link management system - create and manage custom shortlinks on your webserver</strong><br><br>
     <img src="https://forthebadge.com/images/badges/made-with-go.svg" height="30" />&nbsp;
     <img src="https://forthebadge.com/images/badges/made-with-vue.svg" height="30" />&nbsp;
     <img src="https://forthebadge.com/images/badges/60-percent-of-the-time-works-every-time.svg" height="30" />&nbsp;
     <a href="https://zekro.de/discord"><img src="https://img.shields.io/discord/307084334198816769.svg?logo=discord&style=for-the-badge" height="30"></a>
</div>


---

# Introduction

SLMS is a simple server solution which provides managable short links for your domain. You can create, edit and delete short links in a simple web interface which will relocate the user of these links to the defined root link location. Each *(uncached)* access will be recorded *(anonymously)* to provide simple access analytics. 

---

# Screenshots

> *State: Commit [`5f7139a`](https://github.com/zekroTJA/slms/commit/5f7139a38c2b737e906e33304aa8e935ba94297a)*

### Management Page

![](https://i.zekro.de/Code_-_Insiders_WwKjWCP6Cm.png)

### Add Short Links

![](https://i.zekro.de/firefox_yd6F9si7ro.png)

## Edit Short Links 

![](https://i.zekro.de/firefox_dWZdEiXMNy.png)

---

# Why v.2.0?

I realy had intrest on resuming developing on this project and I wanted to get deeper into creating REST API's with go. So I thought about experimenting with `fasthttp` and `fasthttp-routing` instead of `net/http` and `gorilla/mux` *(which I have used in v.1.0)*. Also, I wanted to enhance the project layout and the database structure, so I've decided to re-create this whole project.

---

# 3rd Party Dependencies

- [valayala/fasthttp](https://github.com/valyala/fasthttp)
- [qiangxue/fasthttp-routing](https://github.com/qiangxue/fasthttp-routing)
- [op/go-logging](https://github.com/op/go-logging)
- [ghodss/yaml](https://github.com/ghodss/yaml)
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- [go-gem/sessions](https://github.com/go-gem/sessions)

---

© 2019 Ringo Hoffmann (zekro Development)  
Covered by MIT Licence.
