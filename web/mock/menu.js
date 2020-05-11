export default {
	'/api/v1/menu': function (req, res) {
    setTimeout(() => {
      res.json({
        data: [
          {
            id: 1,
            name: '用户管理',
            seq: 1,
          },
          {
            id: 2,
            name: '数据概览',
            seq: 2,
          }
        ]
      });
    }, 300);
  },
  '/api/v1/submenu': function (req, res) {
    setTimeout(() => {
      res.json({
        data: [
          {
            id: 3,
            parent_id: 1,
            name: '用户列表',
            url: '/user/list',
            icon: 'icon-text',
            detail: '测试1',
            seq: 3
          },
          {
            id: 7,
            parent_id: 1,
            name: '角色管理',
            url: '/user/role',
            icon: 'icon-text',
            detail: '测试2',
            seq: 3
          },
          {
            id: 4,
            parent_id: 1,
            name: '权限管理',
            url: '/user/permission',
            icon: 'icon-text',
            detail: '测试3',
            seq: 3
          },
          {
            id: 5,
            parent_id: 2,
            name: '老板看板',
            url: '/menu/menu',
            icon: 'icon-text',
            detail: '测试4',
            seq: 3
          },
          {
            id: 6,
            parent_id: 2,
            name: '财务看板',
            url: '/menu/submenu',
            icon: 'icon-text',
            detail: '测试5',
            seq: 3
          }
        ]
      });
    }, 300);
  },
  'post /api/v1/submenu/add': function (req, res) {
    console.log(req.body);
    const { username, password } = req.body;
    let responseObj;
    responseObj = {
      code: 0,
      message: '',
      data:  null
    };
    setTimeout(() => {
      res.json(responseObj);
    }, 500);
  },
  'post /api/v1/submenu/del': function (req, res) {
    console.log(req.body);
    let responseObj;
    responseObj = {
      code: 0,
      message: '',
      data:  null
    };
    setTimeout(() => {
      res.json(responseObj);
    }, 500);
  },
};