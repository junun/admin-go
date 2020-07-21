import React, {Fragment, Component} from "react";
import { Modal, Table, Tag, Button } from 'antd';
import {connect} from "dva";

const colorsStatus = [
  {id: 0, name: '成功', color: 'green'},
  {id: 1, name: '异常', color: 'red'},
];

@connect(({ loading, schedule }) => {
    return {
      scheduleHisList: schedule.scheduleHisList,
      scheduleHisLoading: loading.effects['schedule/getScheduleHis'],
    }
})

class Record extends React.Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({
      type: 'schedule/getScheduleHis',
      payload: {
        id: this.props.info.id,
        pagesize: 100, 
      }
    });
  }

  columns = [{
    title: '执行时间',
    dataIndex: 'RunTime'
  }, {
    title: '执行状态',
    dataIndex: 'Status',
    'render': Status => colorsStatus.map(x => {
      if (Status == x.id) {
        return <Tag color={x.color}>{x.name}</Tag>
      }
    })
  }, {
    title: '执行主机',
    dataIndex: 'HostId',
    'render': HostId => this.props.hostList.map(x => {
      if (HostId == x.id) {
        return <Tag>{x.Name}</Tag>
      }
    })
  }, {
    title: '操作',
    render: info => <Button type="link" style={{padding: 0}} onClick={() => this.props.showInfo(2, info.id)}>详情</Button>
  }];

  render() {
    const {scheduleHisList, scheduleHisLoading } = this.props;
    return (
      <Modal
        visible
        width={800}
        maskClosable={false}
        title={`任务执行记录 - ${this.props.info.Name}`}
        onCancel={this.props.onCancel}
        footer={null}>
        <Table columns={this.columns} dataSource={scheduleHisList} loading={scheduleHisLoading} rowKey="id"/>
      </Modal>
    )
  }
}

export default Record