import { getMenus, getNotify, patchNotify} from '@/services/user';
import router from 'umi/router'
import { message } from 'antd';

export default {
  namespace: 'app',

  state: {
    menu: [],
    user: {},
    notifies: [],
  },
  effects: {
    *getMenu(payload, { call, put, select }) {
      var id
      if (sessionStorage.getItem('is_supper')==1) {
        id = 0;
      } else {
        var temp = sessionStorage.getItem('user')
        id = JSON.parse(temp).rid;
      }

      const response = yield call(getMenus, id);
      yield put({
        type: 'updateMenu',
        payload: response.data.lists,
      });
    },
    *getNotify(payload, { call, put, select }){
      const response = yield call(getNotify);
      yield put({
        type: 'updateNotify',
        payload: response.data.lists,
      });
    },
    *patchNotify({ payload }, { call, put, select }){
      const response = yield call(patchNotify, payload);
      if (response && response.code == 200) {
        // yield put({
        //   type: 'getNotify',
        // });
      } else {
        message.error(response.message);
      }
    },
  },
  reducers: {
    updateMenu(state, { payload: menu }) {
      return {
        ...state,
        menu,
      };
    },
    updateUser(state, { payload }) {
      return {
        ...state,
        user: payload,
      }
    },
    updateNotify(state, { payload: notifies }) {
      return {
        ...state,
        notifies,
      };
    },
  },
  subscriptions: {
    setup({ dispatch, history }) {
      history.listen(location => {
        const userJSON = sessionStorage.getItem('user');
        if (userJSON) {
          const user = JSON.parse(userJSON);
          dispatch({
            type: 'updateUser',
            payload: user,
          });
        }
        const pathname = location.pathname;
        if (pathname !== '/user/login') {
          const token = sessionStorage.getItem('jwt');
          const userJSON = sessionStorage.getItem('user');
          if (!token || !userJSON) {
            sessionStorage.removeItem('jwt');
            // 未登录访问
            router.push('/user/login');
          }
        }
      });
    },
  }
};
