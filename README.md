<div align="center">
     <!-- <img src="https://zekro.de/src/go_chat_logo.png" width="400"/> -->
     <h1>~ SLMS ~</h1>
     <strong>Short link management system - create and manage custom shortlinks on your webserver</strong><br><br>
     <img src="https://forthebadge.com/images/badges/made-with-go.svg" height="30" />&nbsp;
     <img src="https://forthebadge.com/images/badges/60-percent-of-the-time-works-every-time.svg" height="30" />&nbsp;
     <a href="https://zekro.de/discord"><img src="https://img.shields.io/discord/307084334198816769.svg?logo=discord&style=for-the-badge" height="30"></a>
</div>


---

# Introduction

SLMS is a server component where you can create and manage shortlinks which redirekt to the set root urls. You can create shortlinks with defined names or random strings. There is a small analytics component counting accesses on your shortlinks and when the last access was taken. All data is saved in a MySql-Database and can be accessed over a web interface.

---

# Screenshots

### Management Page

![](https://cdn.zekro.de/ss/chrome_2018-10-12_11-16-59.png)
*State: Commit `d774cf3ed11e8d8c51c949e5f4b69bae2ad03f03`*

---

# 3rd Party Dependencies

- [gorilla/mux](https://github.com/gorilla/mux)
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- [go-yaml/yaml.v2](https://github.com/go-yaml/yaml/tree/v2.2.1)

---

© 2018 zekro Development  

[zekro.de](https://zekro.de) | contact[at]zekro.de


