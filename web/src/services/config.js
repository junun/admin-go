import { httpGet, httpPost, httpPut, httpDel } from '@/utils/request';
import { stringify } from 'qs';

export async function getConfigEnv(params) {
  return httpGet(`/admin/config/env?${stringify(params)}`);
}

export async function configEnvAdd(params) {
  return httpPost('/admin/config/env', params);
}

export async function configEnvEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/config/env/${id}`, params);
}

export async function configEnvDel(params) {
  return httpDel(`/admin/config/env/${params}`);
}

export async function getAppType(params) {
  return httpGet(`/admin/config/type?${stringify(params)}`);
}

export async function appTypeAdd(params) {
  return httpPost('/admin/config/type', params);
}

export async function appTypeEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/config/type/${id}`, params);
}

export async function appTypeDel(params) {
  return httpDel(`/admin/config/type/${params}`);
}

export async function getProject(params) {
  return httpGet(`/admin/config/app?${stringify(params)}`);
}

export async function configProjectAdd(params) {
  return httpPost('/admin/config/app', params);
}

export async function configProjectEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/config/app/${id}`, params);
}

export async function configProjectDel(params) {
  return httpDel(`/admin/config/app/${params}`);
}

export async function configProjectSync(params) {
  return httpGet(`/admin/config/sync/app/${params}`);
}

export async function getAppValue(params) {
  return httpGet(`/admin/config/value?${stringify(params)}`);
}


export async function appValueAdd(params) {
  return httpPost('/admin/config/value', params);
}

export async function appValueEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/config/value/${id}`, params);
}

export async function appValueDel(params) {
  return httpDel(`/admin/config/value/${params}`);
}

export async function getDeployExtend(params) {
  return httpGet(`/admin/config/deploy?${stringify(params)}`);
}

export async function deployExtendAdd(params) {
  return httpPost('/admin/config/deploy', params);
}

export async function deployExtendEdit(params) {
  var id = params.Dtid
  delete params.Dtid
  return httpPut(`/admin/config/deploy/${id}`, params);
}

export async function deployExtendDel(params) {
  return httpDel(`/admin/config/deploy/${params}`);
}

export async function getAppTemplate(params) {
  return httpGet(`/admin/config/template?${stringify(params)}`);
}
