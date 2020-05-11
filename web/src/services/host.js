import { httpGet, httpPost, httpPut, httpDel } from '@/utils/request';
import { stringify } from 'qs';

export async function getHostRole(params) {
  return httpGet(`/admin/host/role?${stringify(params)}`);
}

export async function hostRoleAdd(params) {
  return httpPost('/admin/host/role', params);
}

export async function hostRoleEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/host/role/${id}`, params);
}

export async function hostRoleDel(params) {
  return httpDel(`/admin/host/role/${params}`);
}

export async function getHost(params) {
  return httpGet(`/admin/host?${stringify(params)}`);
}

export async function hostAdd(params) {
  return httpPost('/admin/host', params);
}

export async function hostEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/hosts/${id}`, params);
}

export async function hostDel(params) {
  return httpDel(`/admin/hosts/${params}`);
}


export async function getHostApp(params) {
  return httpGet(`/admin/host/app?${stringify(params)}`);
}

export async function hostAppAdd(params) {
  return httpPost('/admin/host/app', params);
}

export async function hostAppEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/host/app/${id}`, params);
}

export async function hostAppDel(params) {
  return httpDel(`/admin/host/app/${params}`);
}

export async function getHostByAppId(params) {
  return httpGet(`/admin/host/appid?${stringify(params)}`);
}


