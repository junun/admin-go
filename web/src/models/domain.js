import {getDomain, domainAdd, domainEdit, domainDel} from '@/services/domain';
import router from 'umi/router';
import { message } from 'antd';

export default {
  namespace: 'domain',

  state: {
    domainList: [],
    domainListLen: 0,
    domainPage: 1,
    domainSize: 10,
  },

  reducers: {
    updateDomainPage(state, { payload }) {
      return {
        ...state,
        domainPage: payload.page,
        domainSize: payload.pageSize && payload.pageSize || 100,
      }
    },
    updateDomainList(state, { payload }){
      return {
        ...state,
        domainList: payload.lists,
        domainListLen: payload.count,
      }
    },
  },
  effects: {
    *getDomain({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updateDomainPage',
          payload: payload,
        });
      }
      const state = yield select(state => state.domain);
      const {domainPage, domainSize} = state;
      const query = {
        page: domainPage,
        pagesize: domainSize,
      };
      const response = yield call(getDomain, query);
      yield put({
        type: 'updateDomainList',
        payload: response.data,
      });
    },
    *domainAdd({ payload }, { call, put }) {
      const response = yield call(domainAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getDomain',
        });
      } else {
        message.error(response.message);
      }
    },
    *domainEdit({ payload }, { call, put }) {
      const response = yield call(domainEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getDomain',
        });
      } else {
        message.error(response.message);
      }
    },
    *domainDel({ payload }, { call, put }) {
      const response = yield call(domainDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getDomain',
        });
      } else {
        message.error(response.message);
      }
    },
  }
};
