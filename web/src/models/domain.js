import {getDomain, domainAdd, domainEdit, domainDel, 
  getCertificate, certificateAdd, certificateEdit, 
  certificateDel} from '@/services/domain';
import router from 'umi/router';
import { message } from 'antd';

export default {
  namespace: 'domain',

  state: {
    domainList: [],
    domainListLen: 0,
    domainPage: 1,
    domainSize: 10,

    certificateList: [],
    certificateListLen: 0,
    certificatePage: 1,
    certificateSize: 10,
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
    updateCertificatePage(state, { payload }) {
      return {
        ...state,
        certificatePage: payload.page,
        certificateSize: payload.pageSize && payload.pageSize || 100,
      }
    },
    updateCertificateList(state, { payload }){
      return {
        ...state,
        certificateList: payload.lists,
        certificateListLen: payload.count,
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
    *getCertificate({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updateCertificatePage',
          payload: payload,
        });
      }
      const state = yield select(state => state.domain);
      const {certificatePage, certificateSize} = state;
      const query = {
        page: certificatePage,
        pagesize: certificateSize,
      };
      const response = yield call(getCertificate, query);
      yield put({
        type: 'updateCertificateList',
        payload: response.data,
      });
    },
    *certificateAdd({ payload }, { call, put }) {
      const response = yield call(certificateAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getCertificate',
        });
      } else {
        message.error(response.message);
      }
    },
    *certificateEdit({ payload }, { call, put }) {
      const response = yield call(certificateEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getCertificate',
        });
      } else {
        message.error(response.message);
      }
    },
    *certificateDel({ payload }, { call, put }) {
      const response = yield call(certificateDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getCertificate',
        });
      } else {
        message.error(response.message);
      }
    },
  }
};
