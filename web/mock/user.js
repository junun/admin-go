export default {
  'post /api/v1/user/login1': function (req, res) {
    const { username, password } = req.body;
    let responseObj;
    if (username === 'admin' && password === 'admin') {
      responseObj = {
        code: 0,
        message: '',
        data: {
          "createTime": 1551955400000,
          "sex": 0,
          "headimgurl": "",
          "mobile": "15361429773",
          "nickname": "郭郭",
          "updateTime": 1552032871000,
          "id": 1,
          "status": 0,
          "username": "guolx",
          token: 'abcd-token',
        }
      };
    } else {
      responseObj = {
        code: -1,
        message: '用户名或密码错误',
        data: null
      };
    }
    setTimeout(() => {
      res.json(responseObj);
    }, 500);
  },
  'post /api/v1/upload': function (req, res) {
    
    setTimeout(() => {
      res.json(responseObj);
    }, 500);
  },
  '/api/v1/permission': function (req, res) {
    setTimeout(() => {
      res.json({
        code: 200,
        data: [
          {
            id: 1,
            name: '用户管理',
            child: [
              {
                id: 101,
                name: '用户列表',
                url: '/user/list',
              },
              {
                id: 102,
                name: '角色列表',
                url: '/user/role',
              },
              {
                id: 103,
                name: '权限列表',
                url: '/user/perm',
              },
            ]
          },
          {
            id: 2,
            name: '菜单管理',
            child: [
              {
                id: 201,
                name: '一级菜单',
                url: '/menu/menu',
              },
              {
                id: 202,
                name: '二级菜单',
                url: '/menu/submenu',
              },
            ]
          },
          {
            id: 4,
            name: '公司历程',
            child: [
              {
                id: 401,
                name: '总览',
                url: '/history/list',
              }
            ]
          },
          {
            id: 5,
            name: '新闻咨询',
            child: [
              {
                id: 501,
                name: '咨询列表',
                url: '/info/list',
              },
              {
                id: 502,
                name: 'banner',
                url: '/info/banner',
              },
            ]
          },
          {
            id: 6,
            name: '产品信息',
            child: [
              {
                id: 601,
                name: '产品列表',
                url: '/product/list',
              }
            ]
          },
        ]
      });
    }, 300);
  },
  '/api/account/users1': function (req, res) {
    setTimeout(() => {
      res.json({
        data: [
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 1,
            mobile: "15361429773",
            nickname: "郭郭",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "guolx",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 2,
            mobile: "15361429773",
            nickname: "呼呼",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "huhu",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 3,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 4,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 5,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 6,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 7,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 8,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 9,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 10,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 11,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 12,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 13,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 14,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 15,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 16,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 17,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 18,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 19,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 20,
            mobile: "15361429773",
            nickname: "管理员",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "admin",
          },
          {
            createTime: 1551955400000,
            headimgurl: "",
            id: 21,
            mobile: "15361429773",
            nickname: "呵呵",
            sex: 0,
            role: "admin",
            status: 0,
            token: "abcd-token",
            updateTime: 1552032871000,
            username: "hehe",
          }
        ]
      });
    }, 300);
  },
  '/api/account/roles1': function (req, res) {
    setTimeout(() => {
      res.json({
        data: [
          {
            id: 1,
            detail: "管理员组",
            name: "admin",
            status: 0,
            createTime: 1551953400000,
            updateTime: 1552032871000,
          },
          {
            id: 2,
            detail: "开发组",
            name: "dev",
            status: 0,
            createTime: 1551954400000,
            updateTime: 1552032871000,
          },
        ]
      });
    }, 300);
  },
  '/api/v1/user/allpermission': function (req, res) {
    setTimeout(() => {
      res.json({
        data: [ 
          {
            id: 1,
            name: '用户管理',
            child: [
              {
                id: 3,
                name: "用户列表",
                child: [
                  {
                    id: 11,
                    name: '添加用户',
                  },
                  {
                    id: 12,
                    name: '删除用户',
                  },
                  {
                    id: 13,
                    name: '编辑用户',
                  },
                  {
                    id: 14,
                    name: '重置密码',
                  },
                ],
              },
              {
                id: 4,
                name: "角色列表",
                child: [
                  {
                    id: 15,
                    name: '添加角色',
                  },
                  {
                    id: 16,
                    name: '删除角色',
                  },
                  {
                    id: 17,
                    name: '编辑角色',
                  },
                  {
                    id: 18,
                    name: '角色权限',
                  },
                ],
              },
              {
                id: 5,
                name: "权限列表",
                child: [
                  {
                    id: 19,
                    name: '添加权限',
                  },
                  {
                    id: 20,
                    name: '删除权限',
                  },
                  {
                    id: 21,
                    name: '编辑权限',
                  },
                ],
              },
            ],
          },
          {
            id: 2,
            name: '菜单管理',
            child: [
              {
                id: 6,
                name: "一级菜单",
                child: [
                  {
                    id: 22,
                    name: '添加一级菜单',
                  },
                  {
                    id: 23,
                    name: '删除一级菜单',
                  },
                  {
                    id: 24,
                    name: '修改一级菜单',
                  },
                ],
              },
              {
                id: 7,
                name: "二级菜单",
                child: [
                  {
                    id: 25,
                    name: '添加二级菜单',
                  },
                  {
                    id: 26,
                    name: '删除二级菜单',
                  },
                  {
                    id: 27,
                    name: '修改二级菜单',
                  },
                ],
              },
            ]
          },
          {
            id: 30,
            name: '数据概览',
            child: [
              {
                id: 31,
                name: "老板看板",
                child: [                  
                ],
              },
            ],
          }, 
        ],
      });
    }, 300);
  },
  '/api/v1/user/userpermission': function (req, res) {
    setTimeout(() => {
      res.json({
        code: 0,
        message: '',
        data: ['40', '41', '42'],
      });
    }, 300);
  },
};
