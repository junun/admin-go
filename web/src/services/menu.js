import { httpGet, httpPost, httpPut, httpDel } from '@/utils/request';
import { stringify } from 'qs';

export async function getMenus(params) {
  return httpGet(`/admin/menus?${stringify(params)}`);
}

export async function menuAdd(params) {
  return httpPost('/admin/menus', params);
}

export async function menuEdit(params) {
  return httpPut(`/admin/menus/${params.id}`, params);
}

export async function menuDel(params) {
  return httpDel(`/admin/menus/${params}`);
}

export async function getSubMenus(params) {
  return httpGet(`/admin/submenus?${stringify(params)}`);
}

export async function subMenuAdd(params) {
  return httpPost('/admin/submenus', params);
}

export async function subMenuEdit(params) {
  return httpPut(`/admin/submenus/${params.id}`, params);
}

export async function subMenuDel(params) {
  return httpDel(`/admin/submenus/${params}`);
}