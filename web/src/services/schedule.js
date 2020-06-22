import { httpGet, httpPost, httpPatch, httpPut, httpDel } from '@/utils/request';
import { stringify } from 'qs';

export async function getSchedule(params) {
  return httpGet(`/admin/schedule?${stringify(params)}`);
}

export async function getScheduleHis(params) {
  var id = params.id
  delete params.id
  return httpGet(`/admin/schedule/${id}?${stringify(params)}`);
}

export async function getScheduleInfo(params) {
  var id = params.id
  delete params.id
  return httpGet(`/admin/schedule/${id}/info?${stringify(params)}`);
}

export async function changeScheduleActive(params) {
  return httpPatch(`/admin/schedule`, params);
}

export async function scheduleAdd(params) {
  return httpPost('/admin/schedule', params);
}

export async function scheduleEdit(params) {
  var id = params.id
  delete params.id
  return httpPut(`/admin/schedule/${id}`, params);
}

export async function scheduleDel(params) {
  return httpDel(`/admin/schedule/${params}`);
}