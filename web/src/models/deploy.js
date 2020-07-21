import {getDeploy, deployAdd, deployEdit, deployDel, 
  deployReview, getGitTag, rollbackConfirm,
  getGitBranch, getGitCommit } from '@/services/deploy';
import router from 'umi/router';
import { message } from 'antd';

export default {
  namespace: 'deploy',

  state: {
    deployList: [],
    deployLen: 0,
    deployPage: 1,
    deploySize: 10,
    gitBranchList: [],
    gitCommitList: [],
    gitTagList: [],
  },

  reducers: {
    updatedeployPage(state, { payload }) {
      return {
        ...state,
        deployPage: payload.page,
        deploySize: payload.pageSize && payload.pageSize || 100
      }
    },
    updatedeployList(state, { payload }){
      return {
        ...state,
        deployList: payload.lists,
        deployLen: payload.count,
      }
    },
    updateGitBranch(state, { payload }) {
      return {
        ...state,
        gitBranchList: payload.lists,
      }
    },
    updateGitTag(state, { payload }) {
      return {
        ...state,
        gitTagList: payload.lists,
      }
    },
    updateGitCommit(state, { payload }) {
      return {
        ...state,
        gitCommitList: payload.lists,
      }
    },
    cleanGitInfoList(state, { payload }) {
      return {
        ...state,
        gitBranchList: [],
        gitCommitList: [],
        gitTagList: [],
      }
    },
  },
  effects: {
    *getDeploy({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updatedeployPage',
          payload: payload,
        });
      }
      const state = yield select(state => state.deploy);
      const {deployPage, deploySize} = state;
      const query = {
        page: deployPage,
        pagesize: deploySize,
      };
      const response = yield call(getDeploy, query);
      yield put({
        type: 'updatedeployList',
        payload: response.data,
      });
    },
    *deployAdd({ payload }, { call, put }) {
      const response = yield call(deployAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getDeploy',
        });
      } else {
        message.error(response.message);
      }
    },
    *deployEdit({ payload }, { call, put }) {
      const response = yield call(deployEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getDeploy',
        });
      } else {
        message.error(response.message);
      }
    },
    *deployDel({ payload }, { call, put }) {
      const response = yield call(deployDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getDeploy',
        });
      } else {
        message.error(response.message);
      }
    },
    *deployReview({ payload }, { call, put }) {
      const response = yield call(deployReview, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getDeploy',
        });
        if (payload.IsPass == 1) {
          message.success(response.message);
        } else {
          message.warn(response.message);
        }
      } else {
        message.error(response.message);
      }
    },
    *getGitTag({payload}, { call, put, select }){
      const response = yield call(getGitTag, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'updateGitTag',
          payload: response.data,
        });
      } else {
        response.data.lists = [];
        yield put({
          type: 'updateGitTag',
          payload: response.data,
        });
      }
    },
    *getGitBranch({payload}, { call, put, select }){
      const response = yield call(getGitBranch, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'updateGitBranch',
          payload: response.data,
        });
      } else {
        response.data.lists = [];
        yield put({
          type: 'updateGitBranch',
          payload: response.data,
        });
      }
    },
    *getGitCommit({payload}, { call, put, select }){
      const response = yield call(getGitCommit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'updateGitCommit',
          payload: response.data,
        });
      } else {
        response.data.lists = [];
        yield put({
          type: 'updateGitCommit',
          payload: response.data,
        });
      }
    },
    *cleanBranchList({payload}, { put }){
      yield put({
        type: 'cleanGitInfoList',
      });
    },
    *rollbackConfirm({payload}, { call, put }){
      const response = yield call(rollbackConfirm, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getDeploy',
        });
        message.success(response.message);
      } else {
        message.error(response.message);
      }
    },
    // *getAppVersion({payload}, { call, put, select }){
    //   const response = yield call(getAppVersion, payload);
    //   if (response && response.code == 200) {
    //     console.log(response)
    //   } else {

    //   }
    // },
  }
};
