import {getSchedule, changeScheduleActive, scheduleAdd,
 scheduleEdit, scheduleDel, getScheduleHis, getScheduleInfo} from '@/services/schedule';
import router from 'umi/router';
import { message } from 'antd';

export default {
  namespace: 'schedule',

  state: {
    scheduleList: [],
    scheduleListLen: 0,
    schedulePage: 1,
    scheduleSize: 10,

    scheduleHisList: [],
    scheduleInfo: {},
  },

  reducers: {
    updateSchedulePage(state, { payload }) {
      return {
        ...state,
        schedulePage: payload.page,
        scheduleSize: payload.pageSize && payload.pageSize || 100,
      }
    },
    updateScheduleList(state, { payload }){
      return {
        ...state,
        scheduleList: payload.lists,
        scheduleListLen: payload.count,
      }
    },
    updateScheduleHisList(state, { payload }){
      return {
        ...state,
        scheduleHisList: payload.lists,
      }
    },
    updateScheduleInfo(state, { payload }){
      return {
        ...state,
        scheduleInfo: payload.lists,
      }
    },
  },
  effects: {
    *getSchedule({payload}, {call, put, select }) {
      if (payload) {
        yield put({
          type: 'updateSchedulePage',
          payload: payload,
        });
      }
      const state = yield select(state => state.schedule);
      const {schedulePage, scheduleSize} = state;
      const query = {
        page: schedulePage,
        pagesize: scheduleSize,
      };
      const response = yield call(getSchedule, query);
      yield put({
        type: 'updateScheduleList',
        payload: response.data,
      });
    },
    *getScheduleHis({payload}, {call, put, select }) {
      const response = yield call(getScheduleHis, payload);
      yield put({
        type: 'updateScheduleHisList',
        payload: response.data,
      });
    },
    *getScheduleInfo({payload}, {call, put, select }) {
      const response = yield call(getScheduleInfo, payload);
      yield put({
        type: 'updateScheduleInfo',
        payload: response.data,
      });
    },
    *scheduleAdd({ payload }, { call, put }) {
      const response = yield call(scheduleAdd, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getSchedule',
        });
      } else {
        message.error(response.message);
      }
    },
    *scheduleEdit({ payload }, { call, put }) {
      const response = yield call(scheduleEdit, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getSchedule',
        });
      } else {
        message.error(response.message);
      }
    },
    *scheduleDel({ payload }, { call, put }) {
      const response = yield call(scheduleDel, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getSchedule',
        });
      } else {
        message.error(response.message);
      }
    },
    *changeScheduleActive({ payload }, { call, put }) {
      const response = yield call(changeScheduleActive, payload);
      if (response && response.code == 200) {
        yield put({
          type: 'getSchedule',
        });
      } else {
        message.error(response.message);
      }
    },
  }
};
