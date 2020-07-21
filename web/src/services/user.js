import {httpGet, httpPut, httpPatch, httpPost, httpDel} from '@/utils/request';
import { stringify } from 'qs';


export async function getNotify(params) {
  return httpGet(`/admin/notify`);
}

export async function patchNotify(params) {
  return httpPatch(`/admin/notify`, params);
}

export async function userLogin(params) {
  return httpPost('/admin/user/login', params);
}

export async function userLogout() {
  return httpPost('/admin/user/logout');
}

export async function getSetting() {
  return httpGet('/admin/system');
}

export async function getSettingAbout() {
  return httpGet('/admin/system/about');
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
  return httpGet(`/admin/roles/${params}/permissions`);
}

export async function rolePermsAdd(params) {
  var id = params.id
  delete params.id
  return httpPost(`/admin/roles/${id}/permissions`, params);
}

export async function roleAppAdd(params) {
  var id = params.id
  delete params.id
  return httpPost(`/admin/roles/${id}/app`, params);
}

export async function roleHostAdd(params) {
  var id = params.id
  delete params.id
  return httpPost(`/admin/roles/${id}/host`, params);
}

export async function getAllEnvApp() {
  return httpGet('/admin/env/app');
}

export async function getAllEnvHost() {
  return httpGet('/admin/env/host');
}

export async function getRoleEnvApp(params) {
  return httpGet(`/admin/roles/${params}/app`);
}

export async function getRoleEnvHost(params) {
  return httpGet(`/admin/roles/${params}/host`);
}

export async function getRobot(params) {
  return httpGet(`/admin/system/robot?${stringify(params)}`);
}

export async function robotAdd(params) {
  return httpPost('/admin/system/robot', params);
}

export async function robotEdit(params) {
  return httpPut(`/admin/system/robot/${params.id}`, params);
}

export async function robotDel(params) {
  return httpDel(`/admin/system/robot/${params}`);
}

export async function settingModify(params) {
  return httpPost(`/admin/system`, params);
}

export async function settingMailTest(params) {
  return httpPost(`/admin/system/mail`, params);
}

export async function robotTest(params) {
  return httpPost(`/admin/system/robot/${params.id}`, params);
}



