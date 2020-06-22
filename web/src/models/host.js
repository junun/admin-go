import {getHostRole, hostRoleAdd, hostRoleEdit, hostRoleDel,
       getHost, hostAdd, hostEdit, hostDel, getHostByAppId,
       getHostApp, hostAppAdd, hostAppDel } from '@/services/host';
import router from 'umi/router';
import { message } from 'antd';

export default {
  namespace: 'host',

  state: {
    hostRoleList: [],
    hostRoleLen: 0,
    infoPage: 1,
    infoSize: 10,

    hostAppList: [],
    hostAppLen: 0,
    hostAppPage: 1,
    hostAppSize: 10,
    hostHid: 0,

    hostList: [],
    hostLen: 0,
    hostPage: 1,
    hostSize: 10,
    Rid: 0,
    Name: '',
    Status: 0,
    Source: '',

    hostListByAppId:[],
    hostListByAppIdLen:0
  },

  reducers: {
    updateInfoPage(state, { payload }) {
      return {
        ...state,
        infoPage: payload.page,
        infoSize: payload.pageSize && payload.pageSize || 100
      }
    },
    updateHostPage(state, { payload }) {
      return {
        ...state,
        hostList: [],
        hostPage: payload.page,
        hostSize: payload.pageSize && payload.pageSize || 10,
        Rid : payload.Rid && payload.Rid || 0,
        Name : payload.Name && payload.Name ||  '',
        Source : payload.Source && payload.Source ||  '',
        Status: payload.Status && payload.Status || 0,
      }
    },
    updateHostList(state, { payload }){
      return {
        ...state,
        hostList: payload.lists,
        hostLen: payload.count,
      }
    },
    updateHostAppList(state, { payload }){
      return {
        ...state,
        hostAppList: payload.lists,
        hostAppLen: payload.count,
      }
    },
    updateHostListByAppId(state, { payload }){
      return {
        ...state,
        hostListByAppId: payload.lists,
        hostListByAppIdLen: payload.count,
      }
    },
    updateHostRoleList(state, { payload }){
      return {
        ...state,
        hostRoleList: payload.lists,
        hostRoleLen: payload.count,
      }
    }
  },
  effects: {
    *getHostRole({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updateInfoPage',
          payload: payload,
        });
      }
      const state = yield select(state => state.host);
      const {infoPage, infoSize} = state;
      const query = {
        page: infoPage,
        pagesize: infoSize,
      };
      const response = yield call(getHostRole, query);
      yield put({
        type: 'updateHostRoleList',
        payload: response.data,
      });
    },
    *hostRoleAdd({ payload }, { call, put }) {
      const response = yield call(hostRoleAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getHostRole',
        });
      } else {
        message.error(response.message);
      }
    },
    *hostRoleEdit({ payload }, { call, put }) {
      const response = yield call(hostRoleEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getHostRole',
        });
      } else {
        message.error(response.message);
      }
    },
    *hostRoleDel({ payload }, { call, put }) {
      const response = yield call(hostRoleDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getHostRole',
        });
      } else {
        message.error(response.message);
      }
    },
    *getHost({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updateHostPage',
          payload: payload,
        });
      }
      const state = yield select(state => state.host);
      const {hostPage, hostSize, Rid, Name, Status, Source} = state;
 
      const query = {
        page: hostPage,
        pagesize: hostSize,
        Rid : Rid,
        Name : Name,
        Source: Source,
        Status: Status
      };

      const response = yield call(getHost, query);
      yield put({
        type: 'updateHostList',
        payload: response.data,
      });
    },
    *hostAdd({ payload }, { call, put }) {
      const response = yield call(hostAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getHost',
        });
      } else {
        message.error(response.message);
      }
    },
    *hostEdit({ payload }, { call, put }) {
      const response = yield call(hostEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getHost',
        });
      } else {
        message.error(response.message);
      }
    },
    *hostDel({ payload }, { call, put }) {
      const response = yield call(hostDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getHost',
        });
      } else {
        message.error(response.message);
      }
    },
    *getHostApp({payload}, {call, put, select }) {
      const state = yield select(state => state.host);
      const {hostAppPage, hostAppSize, hostHid} = state;
      const query = {
        page:  payload.page && payload.page || hostAppPage,
        pagesize: payload.pageSize && payload.pageSize || hostAppSize,
        Hid : payload.hid &&  payload.hid || hostHid,
      };

      const response = yield call(getHostApp, query);
      yield put({
        type: 'updateHostAppList',
        payload: response.data,
      });
    },
    *hostAppAdd({ payload }, { call, put }) {
      const response = yield call(hostAppAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getHostApp',
          payload: {
            page: 1,
            pageSize: 10,
            hid: payload.Hid,
          }
        });
      } else {
        message.error(response.message);
      }
    },
    *hostAppDel({ payload }, { call, put }) {
      const response = yield call(hostAppDel, payload.id);
      if (response && response.code == 200) {
        yield put({
          type: 'getHostApp',
          payload: {
            page: 1,
            pageSize: 10,
            hid: payload.Hid,
          }
        });
      } else {
        message.error(response.message);
      }
    },
    *getHostByAppId({payload}, {call, put, select }) {
      const state = yield select(state => state.host);
      const query = {
        Aid : payload.aid &&  payload.aid || 0,
        EnvId: payload.envid &&  payload.envid || 0,
      };

      const response = yield call(getHostByAppId, query);
      yield put({
        type: 'updateHostListByAppId',
        payload: response.data,
      });
    },
  }
};
