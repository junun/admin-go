import { httpGet, httpPost, httpPut, httpDel } from '@/utils/request';
import { stringify } from 'qs';

export async function getDeploy(params) {
  return httpGet(`/admin/deploy/app?${stringify(params)}`);
}

export async function deployAdd(params) {
  return httpPost('/admin/deploy/app', params);
}

export async function deployEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/deploy/app/${id}`, params);
}

export async function deployDel(params) {
  return httpDel(`/admin/deploy/app/${params}`);
}

export async function deployReview(params) {
  return httpPut(`/admin/deploy/app/${params.ID}/review/${params.status}`, params);
}

export async function getGitBranch(params) {
  return httpGet(`/admin/deploy/app/${params}/branch`);
}

export async function getGitCommit(params) {
  return httpGet(`/admin/deploy/app/${params.aid}/commit/${params.name}`);
}