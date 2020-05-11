import React, { PureComponent} from "react";
import {Card, Empty, Table, DatePicker, Row, Col, Select} from "antd";

const Option = Select.Option;
const { RangePicker } = DatePicker;

class TempleTable extends PureComponent {

  render() {
    const { title, data, loading, totalCount, beginDate, endDate, templeSelectData, loadTempleSimple, dateRangeChanged, templeChanged, pageChanged } = this.props;
    /*
    {
        "fanCount": 4,
        "amount": "330",
        "templeId": 1,
        "city": "温州市",
        "name": "西禅寺",
        "lightingRate": 0.1146,
        "divide": "寺庙50%，第三方10%，知恩40%",
        "lightCount": "1980"
      }
     */
    const columns = [
      {
        'title': '排名',
        'dataIndex': 'name',
        'key': 'name',
      },
      {
        'title': '城市',
        'dataIndex': 'city',
        'key': 'city',
      },
      {
        'title': '供灯总收入',
        'dataIndex': 'amount',
        'key': 'amount',
        'render': amount => amount ? (amount / 100).toFixed(2) : '-',
      },
      {
        'title': '灯位总数',
        'dataIndex': 'lightCount',
        'key': 'lightCount',
      },
      {
        'title': '点灯率',
        'dataIndex': 'lightingRate',
        'key': 'lightingRate',
        // 'render': lightingRate => {
        //   if (lightingRate > 0) {
        //     return <span className={styles.upNumber}>{(lightingRate * 100).toFixed(2)}% <Icon type="arrow-up" /></span>
        //   } else if (lightingRate === 0) {
        //     return <span>{lightingRate}</span>
        //   } else {
        //     return <span className={styles.downNumber}>{(Math.abs(lightingRate) * 100).toFixed(2)}% <Icon type="arrow-down" /></span>
        //   }
        // },
        'render': value => (value * 100).toFixed(2) + '%',
      },
      {
        'title': '粉丝量',
        'dataIndex': 'fanCount',
        'key': 'fanCount',
      },
      // {
      //   'title': '目标完成度',
      //   'dataIndex': 'progress',
      //   'key': 'progress',
      //   'render': progress => (
      //     <Progress percent={progress * 100} size="small" />
      //   ),
      // },
      {
        'title': '分成比率',
        'dataIndex': 'divide',
        'key': 'divide',
      },
    ];
    const picker = <RangePicker value={[beginDate, endDate]} onChange={dateRangeChanged} />;
    const select = <Select
      disabled={loadTempleSimple}
      onChange={templeChanged}
      showSearch
      style={{ width: 200 }}
      optionFilterProp="children"
      defaultValue="0"
      filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
    >
      <Option value="0">全部寺庙</Option>
      {templeSelectData.map(x => <Option key={x.templeId} value={x.templeId}>{x.name}</Option>)}
    </Select>;
    const extra = <Row gutter={16}><Col span={12}>{picker}</Col><Col span={12}>{select}</Col></Row>;

    const expandedRowRender = record => {
      const columns = [
        { title: '价格', dataIndex: 'price', key: 'price',
          'render': x => x ? (x / 100).toFixed(2) : '-',
        },
        { title: '数量', dataIndex: 'quantity', key: 'quantity' },
        { title: '占比', dataIndex: 'rate', key: 'rate',
          'render': x => (x * 100).toFixed(2) + '%',
        },
        { title: '合计', dataIndex: 'amount', key: 'amount',
          'render': x => x ? (x / 100).toFixed(2) : '-', },
      ];
      const data = record.priceRate.sort(function(a, b) {
        return b.price - a.price;
      });
      return <Table
        columns={columns}
        dataSource={data}
        pagination={false}
        size="middle"
      />;
    };
    return <Card title={title} extra={extra}>
      {(data && data.content.length > 0) ?
        <Table
          rowKey="templeId"
          loading={loading}
          dataSource={data.content}
          columns={columns}
          expandRowByClick={true}
          expandedRowRender={expandedRowRender}
          pagination={{
          defaultCurrent:1,
          total: totalCount,
          showQuickJumper: true,
          showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
          onChange: pageChanged}} />
      :
        <Empty />
      }
    </Card>
  }
}

export default TempleTable;
