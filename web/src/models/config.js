import {getConfigEnv, configEnvAdd, configEnvEdit, configEnvDel, getProject,
  configProjectAdd, configProjectEdit, configProjectSync, configProjectDel,
  appTypeAdd, appTypeEdit, appTypeDel, getAppType, getAppTemplate,
  getDeployExtend, deployExtendAdd, deployExtendEdit, deployExtendDel,
  getAppValue, appValueAdd, appValueEdit, appValueDel } from '@/services/config';
import router from 'umi/router';
import { message } from 'antd';

export default {
  namespace: 'config',

  state: {
    deployTempleList: [],
    configEnvList: [],
    configEnvLen: 0,
    envPage: 1,
    envSize: 10,
    appTypeList: [],
    appTypeLen: 0,
    appTypePage: 1,
    appTypeSize: 10,
    // 项目列表
    projectsList: [],
    // 项目列表
    appTemplateList: [],
    // 总项目数
    projectsLen: 0,
    // 项目列表当前页
    projectPage: 1,
    // 项目列表分页大小
    projectSize: 10,
    active: 0,
    envId: 0,
    appValueList: [],
    appValueLen: 0,
    appValuePage: 1,
    appValueSize: 10,
    appAid: 0,
    data: ["###starting###"],
  },

  reducers: {
    updateProjectPage(state, { payload }) {
      return {
        ...state,
        projectPage: payload.page,
        projectSize: payload.pageSize && payload.pageSize || 100,
        active: payload.active &&  payload.active  || 0,
        envId: payload.envId &&  payload.envId || 0,
      }
    },
    updateProjectsList(state, { payload }){
      return {
        ...state,
        projectsList: payload.lists,
        projectsLen: payload.count,
      }
    },
    updateAppTemplateList(state, { payload }){
      return {
        ...state,
        appTemplateList: payload.lists,
      }
    },
    updateEnvPage(state, { payload }) {
      return {
        ...state,
        envPage: payload.page,
        envSize: payload.pageSize && payload.pageSize || 100
      }
    },
    updateAppTypePage(state, { payload }) {
      return {
        ...state,
        appTypePage: payload.page,
        appTypeSize: payload.pageSize && payload.pageSize || 100
      }
    },
    updateAppValuePage(state, { payload }) {
      return {
        ...state,
        appValuePage: payload.page,
        appValueSize: payload.pageSize && payload.pageSize || 10,
        appAid: payload.aid && payload.aid || 0,
      }
    },
    updateAppValue(state, { payload }) {
      return {
        ...state,
        appValueList: payload.lists,
        appValueLen: payload.count,
      }
    },
    updateConfigEnvList(state, { payload }){
      return {
        ...state,
        configEnvList: payload.lists,
        configEnvLen: payload.count,
      }
    },
    updateAppTypeList(state, { payload }){
      return {
        ...state,
        appTypeList: payload.lists,
        appTypeLen: payload.count,
      }
    },
    updateData(state, { payload }){
      return {
        ...state,
        data: payload,
      }
    },
    updateDeployExtendList(state, { payload }){
      return {
        ...state,
        deployTempleList: payload.lists,
      }
    },
  },
  effects: {
    *getConfigEnv({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updateEnvPage',
          payload: payload,
        });
      }
      const state = yield select(state => state.config);
      const {envPage, envSize} = state;
      const query = {
        page: envPage,
        pagesize: envSize,
      };
      const response = yield call(getConfigEnv, query);
      yield put({
        type: 'updateConfigEnvList',
        payload: response.data,
      });
    },
    *configEnvAdd({ payload }, { call, put }) {
      const response = yield call(configEnvAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getConfigEnv',
        });
      } else {
        message.error(response.message);
      }
    },
    *configEnvEdit({ payload }, { call, put }) {
      const response = yield call(configEnvEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getConfigEnv',
        });
      } else {
        message.error(response.message);
      }
    },
    *configEnvDel({ payload }, { call, put }) {
      const response = yield call(configEnvDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getConfigEnv',
        });
      } else {
        message.error(response.message);
      }
    },
    *getAppType({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updateAppTypePage',
          payload: payload,
        });
      }
      const state = yield select(state => state.config);
      const {appTypePage, appTypeSize} = state;
      const query = {
        page: appTypePage,
        pagesize: appTypeSize,
      };
      const response = yield call(getAppType, query);
      yield put({
        type: 'updateAppTypeList',
        payload: response.data,
      });
    },
    *appTypeAdd({ payload }, { call, put }) {
      const response = yield call(appTypeAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getAppType',
        });
      } else {
        message.error(response.message);
      }
    },
    *appTypeEdit({ payload }, { call, put }) {
      const response = yield call(appTypeEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getAppType',
        });
      } else {
        message.error(response.message);
      }
    },
    *appTypeDel({ payload }, { call, put }) {
      const response = yield call(appTypeDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getAppType',
        });
      } else {
        message.error(response.message);
      }
    },
    *getAppValue({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updateAppValuePage',
          payload: payload,
        });
      }
      const state = yield select(state => state.config);
      const { appValueSize, appValuePage, appAid  } = state;

      const query = {
        page: appValuePage,
        pageSize: appValueSize,
        Aid: appAid,
      };

      const response = yield call(getAppValue, query);
      yield put({
        type: 'updateAppValue',
        payload: response.data,
      });
    },
    *appValueAdd({ payload }, { call, put }) {
      const response = yield call(appValueAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getAppValue',
          payload: {
            page: 1,
            pageSize: 10,
            aid: payload.Aid,
          }
        });
      } else {
        message.error(response.message);
      }
    },
    *appValueEdit({ payload }, { call, put }) {
      const response = yield call(appValueEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getAppValue',
          payload: {
            page: 1,
            pageSize: 10,
            aid: payload.Aid,
          }
        });
      } else {
        message.error(response.message);
      }
    },
    *appValueDel({ payload }, { call, put }) {
      const response = yield call(appValueDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getAppValue',
          payload: {
            page: 1,
            pageSize: 10,
            aid: payload.Aid,
          }
        });
      } else {
        this
        message.error(response.message);
      }
    },
    *getProject({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updateProjectPage',
          payload: payload,
        });
      }
      const state = yield select(state => state.config);
      const { projectPage, projectSize , active, envId } = state;

      const query = {
        page: projectPage,
        pageSize: projectSize,
        active:  active,
        envId: envId,
      };

      const response = yield call(getProject, query);
      yield put({
        type: 'updateProjectsList',
        payload: response.data,
      });
    },
    *configProjectAdd({ payload }, { call, put }) {
      const response = yield call(configProjectAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getProject',
        });
      } else {
        message.error(response.message);
      }
    },
    *configProjectEdit({ payload }, { call, put }) {
      const response = yield call(configProjectEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getProject',
        });
      } else {
        message.error(response.message);
      }
    },
    *configProjectDel({ payload }, { call, put }) {
      const response = yield call(configProjectDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getProject',
        });
      } else {
        this
        message.error(response.message);
      }
    },
    *configProjectSync({ payload }, { call, put, select }) {
      const response = yield call(configProjectSync, payload);
      const state = yield select(state => state.config);
      const { data} = state;
      data.push(response.message)

      yield put({
        type: 'updateData',
        payload: data,
      });
      if (response && response.code == 200) {
        message.success(response.message);
      } else {
        message.error(response.message);
      }
    },
    *getDeployExtend({payload}, {call, put, select }) {
      const query = {
        Aid:  payload.id,
      };
      const response = yield call(getDeployExtend, query);
      yield put({
        type: 'updateDeployExtendList',
        payload: response.data,
      });
    },
    *deployExtendAdd({ payload }, { call, put }) {
      const response = yield call(deployExtendAdd, payload);
      if (response && response.code == 200) {
        window.location.reload();
      } else {
        message.error(response.message);
      }
    },
    *deployExtendEdit({ payload }, { call, put }) {
      console.log(payload);
      const response = yield call(deployExtendEdit, payload);
      if (response && response.code == 200) {
        window.location.reload();
      } else {
        message.error(response.message);
      }
    },
    *deployExtendDel({ payload }, { call, put }) {
      const response = yield call(deployExtendDel, payload.tid);
      if (response && response.code == 200) {
        window.location.reload();
      } else {
        message.error(response.message);
      }
    },
    *getAppTemplate({payload}, {call, put, select }) {
      const query = {
        aid: payload.aid,
      };
      const response = yield call(getAppTemplate, query);
      yield put({
        type: 'updateAppTemplateList',
        payload: response.data,
      });
    },
  }
};
