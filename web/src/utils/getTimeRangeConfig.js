import moment from "moment";

// 时间范围配置

const incomeConfig = {
  '0': { // 今日
    begin: () => {
      return moment().startOf('day');
    },
    end: () => {
      return moment().endOf('day');
    },
    type: 'HOUR',
  },
  '1': { // 周
    begin: () => {
      return moment().subtract(7, 'days').startOf('day');
    },
    end: () => {
      return moment().endOf('day');
    },
    type: 'DAY',
  },
  '2': { // 月
    begin: () => {
      return moment().subtract(1, 'month').startOf('day');
    },
    end: () => {
      return moment().endOf('day');
    },
    type: 'DAY',
  },
  '3': { // 季度
    begin: () => {
      return moment().subtract(3, 'month').startOf('day');
    },
    end: () => {
      return moment().endOf('day');
    },
    type: 'DAY',
  },
  '4': { // 年
    begin: () => {
      return moment().subtract(1, 'year').startOf('day');
    },
    end: () => {
      return moment().endOf('day');
    },
    type: 'MONTH',
  },
};

const getTimeRangeConfig = (config) => {
  return incomeConfig[config];
};

export default getTimeRangeConfig;
