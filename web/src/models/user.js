import { userLogin, userLogout, getLists, getRoles, getPermissions, permAdd, 
  permEdit, getAllPermissions, getUserPermissions, userAdd, userEdit, 
  settingModify, settingMailTest, robotTest, getRoleEnvApp,
  getAllEnvApp, roleAppAdd, getAllEnvHost, getRoleEnvHost, roleHostAdd,
  getSetting, getSettingAbout, getRobot, robotAdd, robotEdit, robotDel,
  userDel, roleAdd, roleEdit, roleDel, permDel, rolePermsAdd, getRolePerms } from '@/services/user';
import router from 'umi/router';
import { message } from 'antd';

export default {
  namespace: 'user',

  state: {
    //用户列表
    usersList: [],
    usersCount: 0,
    rolesList: [],
    rolesCount: 0,
    permissionsList : [],
    permissionsTotal: 0,
    allPermissionsList : [],
    userPermissionsList : [],
    roleVisible: false,
    // 权限列表当前页
    permPage: 1,
    // 权限列表分页大小
    permSize: 10,
    settingList: [],
    settingAbout: {}, 
    allEnvAppList: [],
    envAppList: {},
    allEnvHostList: [],
    envHostList: {},
    robotList: [],
    robotCount: 0,
    robotPage: 1,
    robotSize: 10,
  },

  reducers: {
    updateList(state, { payload }) {
      return {
        ...state,
        usersList: payload.lists,
        usersCount: payload.count,
      }
    },
    updateRoleList(state, { payload }) {
      return {
        ...state,
        rolesList: payload.lists,
        rolesCount: payload.count,
      }
    },
    updatePermissionList(state, { payload }){
      return {
        ...state,
        permissionsList: payload.lists,
        permissionsTotal: payload.count,
      }
    },
    updateAllPermissionList(state, { payload }){
      return {
        ...state,
        allPermissionsList: payload,
      }
    },
    updateAllEnvAppList(state, { payload }){
      return {
        ...state,
        allEnvAppList: payload,
      }
    },
    updateAllEnvHostList(state, { payload }){
      return {
        ...state,
        allEnvHostList: payload,
      }
    },
    updateRoleEnvAppList(state, { payload }){
      return {
        ...state,
        envAppList: payload,
      }
    },
    updateRoleEnvHostList(state, { payload }){
      return {
        ...state,
        envHostList: payload,
      }
    },
    updateUserPermissionList(state, { payload }){
      return {
        ...state,
        userPermissionsList: payload,
        // roleVisible: true,
      }
    },
    cancelUserPermission(state) {
      return {
        ...state,
        roleVisible: false,
        userPermissionsList: [],
      }
    },
    updatePermPage(state, { payload }) {
      return {
        ...state,
        permPage: payload.page && payload.page || 1,
        permSize: payload.hasOwnProperty('pagesize') && payload.pagesize || 10,
      }
    },
    updateSettingList(state, { payload }) {
      return {
        ...state,
        settingList: payload.lists,
      }
    },
    updateSettingAbout(state, { payload }) {
      return {
        ...state,
        settingAbout: payload.lists,
      }
    },
    updateRobotList(state, { payload }) {
      return {
        ...state,
        robotList: payload.lists,
        robotCount: payload.count,
      }
    },
  },

  effects: {
    *login({ payload }, { call, put }) {
      const response = yield call(userLogin, payload);
      yield put({
        type: 'changeLoginStatus',
        payload: response,
      });
      if (response && response.code == 200) {
        const token = response.data.token;
        sessionStorage.setItem('jwt', token);
        sessionStorage.setItem('is_supper', response.data.is_supper);
        sessionStorage.setItem('permissions', response.data.permissions);
        sessionStorage.setItem('user', JSON.stringify(response.data));
        yield put(router.push('/welcome'));
      } else {
        message.error(response.message);
      }
    },
    *logout(payload, { call, put }) {
      const response = yield call(userLogout);
      if (response && response.code == 200) {
        sessionStorage.removeItem('jwt');
        sessionStorage.removeItem('user');
        sessionStorage.removeItem('is_supper');
        sessionStorage.removeItem('permissions');
        yield put(router.push('/user/login'))
      }
    },
    *getRole({payload}, { call, put, select }){
      const response = yield call(getRoles);
      yield put({
        type: 'updateRoleList',
        payload: response.data,
      });
    },
    *getPermission({payload}, { call, put, select }) {
      if (payload) {
        yield put({
          type: 'updatePermPage',
          payload: payload,
        });
      }
      const state = yield select(state => state.user);
      const { permPage, permSize } = state;
      const query = {
        page: permPage,
        pagesize: permSize,
      };
    
      const response = yield call(getPermissions, query);
      yield put({
        type: 'updatePermissionList',
        payload: response.data,
      });
    },

    *permAdd({payload}, { call, put, select }){
      const response = yield call(permAdd, payload);
      if (response && response.code == 200) {
        // 更新角色列表
        yield put({
          type: 'getPermission',
        });
      } else {
        message.error(response.message);
      }
    },
    *permEdit({payload}, { call, put, select }){
      const response = yield call(permEdit, payload);
      if (response && response.code == 200) {
        // 更新权限列表
        yield put({
          type: 'getPermission',
        });
      } else {
        message.error(response.message);
      }
    },
    *permDel({payload}, { call, put, select }){
      const response = yield call(permDel, payload);
      if (response && response.code == 200) {
        // 更新角色列表
        yield put({
          type: 'getPermission',
        });
      } else {
        message.error('删除错误， 请检查是否有依赖？');
      }
    },
    *getList({payload}, { call, put, select }){
      const response = yield call(getLists);
      yield put({
        type: 'updateList',
        payload: response.data,
      });
    },
    *userAdd({payload}, { call, put, select }){
      const response = yield call(userAdd, payload);
      if (response && response.code == 200) {
        // 更新角色列表
        yield put({
          type: 'getList',
        });
      } else {
        message.error(response.message);
      }
    },
    *userEdit({payload}, { call, put, select }){
      const response = yield call(userEdit, payload);
      if (response && response.code == 200) {
        // 更新角色列表
        yield put({
          type: 'getList',
        });
      } else {
        message.error(response.message);
      }
    },
    *userDel({payload}, { call, put, select }){
      const response = yield call(userDel, payload);
      if (response && response.code == 200) {
        // 更新角色列表
        yield put({
          type: 'getList',
        });
      } else {
        message.error(response.message);
      }
    },
    *roleAdd({payload}, { call, put, select }){
      const response = yield call(roleAdd, payload);
      if (response && response.code == 200) {
        // 更新角色列表
        yield put({
          type: 'getRole',
        });
      } else {
        message.error(response.message);
      }
    },
    *roleEdit({payload}, { call, put, select }){
      const response = yield call(roleEdit, payload);
      if (response && response.code == 200) {
        // 更新角色列表
        yield put({
          type: 'getRole',
        });
      } else {
        message.error(response.message);
      }
    },
    *roleDel({payload}, { call, put, select }){
      const response = yield call(roleDel, payload);
      // response = JSON.paser(response);
      if (response && response.code == 200) {
        // 更新角色列表
        yield put({
          type: 'getRole',
        });
      } else {
        message.error('删除错误， 请检查是否有依赖？');
      }
    },
    *getAllPermission({payload}, { call, put, select }){
      const response = yield call(getAllPermissions);
      if (response && response.message === '' ) {
        yield put({
          type: 'updateAllPermissionList',
          payload: response.data.lists,
        });
      }
    },
    *getRolePerms({payload}, { call, put, select }){
      const response = yield call(getRolePerms, payload);
      const temp = [];

      response.data.lists.map(item => {
        // temp.push((item.Pid).toString());
        temp.push((item.Pid));
      });

      yield put({
        type: 'updateUserPermissionList',
        payload: temp,
      });
    },
    *getAllEnvApp({payload}, { call, put, select }){
      const response = yield call(getAllEnvApp);
      if (response && response.code == 200) {
        yield put({
          type: 'updateAllEnvAppList',
          payload: response.data,
        });
      }
    },
    *getAllEnvHost({payload}, { call, put, select }){
      const response = yield call(getAllEnvHost);
      if (response && response.code == 200) {
        yield put({
          type: 'updateAllEnvHostList',
          payload: response.data,
        });
      }
    },
    *getRoleEnvApp({payload}, { call, put, select }){
      const response = yield call(getRoleEnvApp, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'updateRoleEnvAppList',
          payload: response.data,
        });
      }
    },
    *getRoleEnvHost({payload}, { call, put, select }){
      const response = yield call(getRoleEnvHost, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'updateRoleEnvHostList',
          payload: response.data,
        });
      }
    },
    *roleAppAdd({payload}, { call, put, select }){
      var id = payload.id
      const response = yield call(roleAppAdd, payload);
      if (response && response.code == 200) {
        message.success(response.message)
      } else {
        message.error(response.message);
      }
    },
    *roleHostAdd({payload}, { call, put, select }){
      var id = payload.id
      const response = yield call(roleHostAdd, payload);
      if (response && response.code == 200) {
        message.success(response.message)
      } else {
        message.error(response.message);
      }
    },
    *rolePermsAdd({payload}, { call, put, select }){
      var id = payload.id
      const response = yield call(rolePermsAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getRolePerms',
          payload: id,
        });
      } else {
        message.error(response.message);
      }
    },
    *getSetting({payload}, { call, put, select }){
      const response = yield call(getSetting);
      yield put({
        type: 'updateSettingList',
        payload: response.data,
      });
    },
    *getSettingAbout({payload}, { call, put, select }){
      const response = yield call(getSettingAbout);
      yield put({
        type: 'updateSettingAbout',
        payload: response.data,
      });
    },
    *settingModify({payload}, { call, put, select }){
      const response = yield call(settingModify, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getSetting',
        });
        message.success(response.message);
      } else {
        message.error(response.message);
      }
    },
    *settingMailTest({payload}, { call, put, select }){
      const response = yield call(settingMailTest, payload);
      if (response && response.code == 200) {
        message.success(response.message);
      } else {
        message.error(response.message);
      }
    },
    *getRobot({payload}, { call, put, select }){
      var query = {};
      if (payload) {
        query = {
          page: payload.page,
          pagesize: payload.pageSize,
        };
        if ('status' in payload) {
          query.status = payload.status
        }
      }
      const response = yield call(getRobot, query);
      yield put({
        type: 'updateRobotList',
        payload: response.data,
      });
    },
    *robotAdd({payload}, { call, put, select }){
      const response = yield call(robotAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getRobot',
        });
      } else {
        message.error(response.message);
      }
    },
    *robotEdit({payload}, { call, put, select }){
      const response = yield call(robotEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getRobot',
        });
      } else {
        message.error(response.message);
      }
    },
    *robotDel({payload}, { call, put, select }){
      const response = yield call(robotDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getRobot',
        });
      } else {
        message.error(response.message);
      }
    },
    *robotTest({payload}, { call, put, select }){
      const response = yield call(robotTest, payload);
      if (response && response.code == 200) {
        message.success(response.message);
      } else {
        message.error(response.message);
      }
    },
  },
};
