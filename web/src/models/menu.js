import { getMenus, getSubMenus, menuAdd, menuEdit, menuDel, subMenuAdd, subMenuEdit, subMenuDel } from '@/services/menu';
import router from 'umi/router';
import { message } from 'antd';

export default {
  namespace: 'menu',

  state: {
    // 菜单列表
    menusList: [],
    subMenusList: [],
    subMenusLen: 0,
    // 二级菜单当前页
    subMenuPage: 1,
    // 二级菜单分页大小
    subMenuSize: 100,
  },
  
  reducers: {
    updateMenuList(state, { payload }){
      return {
        ...state,
        menusList: payload,
      }
    },
    updateSubMenuList(state, { payload }){
      return {
        ...state,
        subMenusList: payload.lists,
        subMenusLen: payload.count,
      }
    },
    updateSubMenuPage(state, { payload }) {
      return {
        ...state,
        type: 1,
        subMenuPage: payload.page,
        subMenuSize: payload.pageSize && payload.pageSize || 100
      }
    },
  },
  effects: {
    *getMenu({payload}, { call, put, select }){
      if (payload) {
        yield put({
          type: 'updateSubMenuPage',
          payload: payload,
        });
      }
      const state = yield select(state => state.menu);
      const { subMenuPage, subMenuSize } = state;
      const query = {
        page: subMenuPage,
        pagesize: subMenuSize,
      };

      const response = yield call(getMenus, query);

      yield put({
        type: 'updateMenuList',
        payload: response.data.lists,
      });
    },
    *menuAdd({ payload }, { call, put }) {
      const response = yield call(menuAdd, payload);
      if (response && response.code == 200) {
        // 更新主菜单
        // yield put({
        //   type: 'getMenu',
        // });
        window.location.reload();
      } else {
        message.error(response.message);
      }
    },
    *menuEdit({ payload }, { call, put }) {
      const response = yield call(menuEdit, payload);
      if (response && response.code == 200) {
        // 更新主菜单
        // yield put({
        //   type: 'getMenu',
        // });
        window.location.reload();
      } else {
        message.error(response.message);
      }
    },
    *menuDel({ payload }, { call, put }) {
      const response = yield call(menuDel, payload);
      if (response && response.code == 200) {
        // 更新主菜单
        // yield put({
        //   type: 'getMenu',
        // });
        window.location.reload();
      } else {
        message.error(response.message);
      }
    },
    *getSubMenu({payload}, { call, put, select }){
      if (payload) {
        yield put({
          type: 'updateSubMenuPage',
          payload: payload,
        });
      }
      const state = yield select(state => state.menu);
      const { subMenuPage, subMenuSize } = state;
      const query = {
        page: subMenuPage,
        pagesize: subMenuSize,
      };
      
      const response = yield call(getSubMenus, query);
      yield put({
        type: 'updateSubMenuList',
        payload: response.data,
      });
    },
    *subMenuAdd({ payload }, { call, put }) {
      const response = yield call(subMenuAdd, payload);
      if (response && response.code == 200) {
        // 更新二级菜单
        // yield put({
        //   type: 'getSubMenu',
        // });
        // // 更新主菜单
        // yield put({
        //   type: 'app/getMenu',
        // });
        window.location.reload();
      } else {
        message.error(response.message);
      }
    },
    *subMenuEdit({ payload }, { call, put }) {
      const response = yield call(subMenuEdit, payload);
      if (response && response.code == 200) {
        // 更新二级菜单
        // yield put({
        //   type: 'getSubMenu',
        // });
        // // 更新主菜单
        // yield put({
        //   type: 'app/getMenu',
        // });
        window.location.reload();
      } else {
        message.error(response.message);
      }
    },
    *subMenuDel({ payload }, { call, put }){
      const response = yield call(subMenuDel, payload);
      if (response && response.code == 200) {
        // 更新二级菜单
        // yield put({
        //   type: 'getSubMenu',
        // });
        window.location.reload();
      } else {
        message.error(response.message);
      }
    },
  },
};
