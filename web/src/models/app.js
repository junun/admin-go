import { getMenus} from '@/services/user';
import router from 'umi/router'

export default {
  namespace: 'app',

  state: {
    menu: [],
    user: {},
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
