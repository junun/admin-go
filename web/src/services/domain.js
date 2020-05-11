import { httpGet, httpPost, httpPut, httpDel } from '@/utils/request';
import { stringify } from 'qs';

export async function getDomain(params) {
  return httpGet(`/admin/domain/info?${stringify(params)}`);
}

export async function domainAdd(params) {
  return httpPost('/admin/domain/info', params);
}

export async function domainEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/domain/info/${id}`, params);
}

export async function domainDel(params) {
  return httpDel(`/admin/domain/info/${params}`);
}

export async function getCertificate(params) {
  return httpGet(`/admin/domain/cert?${stringify(params)}`);
}

export async function certificateAdd(params) {
  return httpPost('/admin/domain/cert', params);
}

export async function certificateEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/domain/cert/${id}`, params);
}

export async function certificateDel(params) {
  return httpDel(`/admin/domain/cert/${params}`);
}