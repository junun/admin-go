import {httpGet, httpPut, httpPost, httpDel} from '@/utils/request';
import { stringify } from 'qs';

export async function userLogin(params) {
  return httpPost('/admin/user/login', params);
}

export async function userLogout() {
  return httpPost('/admin/user/logout');
}

export async function getMenu() {
  return httpGet('/api/v1/permission');
}

export async function getMenus(params) {
  return httpGet(`/admin/user/perms/${params}`);
}

export async function getLists() {
  return httpGet('/admin/user');
}

export async function userAdd(params) {
  return httpPost('/admin/user', params);
}
export async function userEdit(params) {
  return httpPut(`/admin/user/${params.id}`, params);
}

export async function userDel(params) {
  return httpDel(`/admin/user/${params}`);
}

export async function getRoles() {
  return httpGet('/admin/roles');
}

export async function roleAdd(params) {
  return httpPost('/admin/roles', params);
}

export async function roleEdit(params) {
  return httpPut(`/admin/roles/${params.id}`, params);
}

export async function roleDel(params) {
  return httpDel(`/admin/roles/${params.id}`);
}

export async function getPermissions(params) {
  return httpGet(`/admin/perms?${stringify(params)}`);
}

export async function permAdd(params) {
  return httpPost('/admin/perms', params);
}

export async function permEdit(params) {
  return httpPut(`/admin/perms/${params.id}`, params);
}

export async function permDel(params) {
  return httpDel(`/admin/perms/${params}`);
}

export async function getAllPermissions() {
  return httpGet('/admin/perms/lists');
}

export async function getRolePerms(params) {
  // return httpGet(`/api/v1/user/vpns?${stringify(params)}`);
  return httpGet(`/admin/roles/${params}/permissions`);

}

export async function rolePermsAdd(params) {
  var id = params.id
  delete params.id
  return httpPost(`/admin/roles/${id}/permissions`, params);
}

