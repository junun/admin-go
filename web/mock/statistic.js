export default {
  '/api/v1/statistic/summary': function (req, res) {
    setTimeout(() => {
      res.json({
        code: 0,
        data: {
          visit: 124345,
          visitCompare: -0.1,
          temple: 28,
          templeCompare: 0.26,
          light: 680,
          lightCompare: 0.26,
        },
      });
    }, 500);
  },
  '/api/v1/statistic/income': function (req, res) {
    const { begin, end, type } = req.query;
    let data = null;
    if (type === 'hour') {
      data = [
        { date: '00:00', value: 0 },
        { date: '02:00', value: 0 },
        { date: '04:00', value: 0 },
        { date: '06:00', value: 0 },
        { date: '08:00', value: 0 },
        { date: '10:00', value: 666 },
        { date: '12:00', value: 19.9 },
        { date: '14:00', value: 30 },
        { date: '16:00', value: 999 },
        { date: '18:00', value: 33 },
        { date: '20:00', value: 0 },
        { date: '22:00', value: 0 },
        { date: '24:00', value: 0 },
      ];
    } else if (type === 'day') {
      data = [
        { date: '03-06', value: 328 },
        { date: '03-07', value: 122 },
        { date: '03-08', value: 222 },
        { date: '03-09', value: 599 },
        { date: '03-10', value: 1200 },
        { date: '03-11', value: 666 },
        { date: '03-12', value: 999 },
      ];
    } else if (type === 'month') {
      data = [
        { date: '19年1月', value: 328 },
        { date: '19年2月', value: 122 },
        { date: '19年3月', value: 222 },
        { date: '19年4月', value: 599 },
        { date: '19年5月', value: 1200 },
        { date: '19年6月', value: 666 },
        { date: '19年7月', value: 999 },
        { date: '19年8月', value: 328 },
        { date: '19年9月', value: 122 },
        { date: '19年10月', value: 222 },
        { date: '19年11月', value: 599 },
        { date: '19年12月', value: 1200 },
      ];
    }
    setTimeout(() => {
      res.json({
        code: 0,
        data: data,
      });
    }, 500);
  },
};
