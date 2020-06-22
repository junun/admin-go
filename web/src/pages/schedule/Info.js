import React, {Fragment, Component} from "react";
import { Modal, Table, Tag, Spin, Card, Tabs, Form, Row, Col } from 'antd';
import moment from 'moment';
import {isEmpty, timeDatetimeTrans} from "@/utils/globalTools";


class InfoModal extends React.Component {
  constructor(props) {
    super(props);
  }

  showHostName = (id) => {
    for (var i = 0; i < this.props.hostList.length; i++) {
      if (id  == this.props.hostList[i].id) {
        return this.props.hostList[i].Name
      }
    }
  };

  render() {
    const preStyle = {
      marginTop: 5,
      backgroundColor: '#eee',
      borderRadius: 5,
      padding: 10,
      maxHeight: 215,
    };

    return (
      <Modal
        visible
        width={800}
        maskClosable={false}
        title="任务执行详情"
        onCancel={this.props.onCancel}
        footer={null}>
        {
        this.props.scheduleInfo != null &&
          <Spin spinning={this.props.scheduleInfoLoading}>
            <Row gutter={16}>
              <Col span={12}>
                <Card title="执行成功" bordered={false}>
                  {this.props.scheduleInfo.Success}
                </Card>
              </Col>
              <Col span={12}>
                <Card title="执行失败" bordered={false}>
                  {this.props.scheduleInfo.Failure}
                </Card>
              </Col>
            </Row>
            {this.props.scheduleInfo.Outputs && (
              <Tabs tabPosition="left" defaultActiveKey="0" style={{width: 700, height: 350, margin: 'auto'}}>
                {this.props.scheduleInfo.Outputs.map((item, index) => (
                  <Tabs.TabPane
                    key={`${index}`}
                    tab={item.Status == 0 ? this.showHostName(item.HostId) : <span style={{color: 'red'}}>{this.showHostName(item.HostId)}</span>}
                  >
                    <div>执行时间：{timeDatetimeTrans(item.CreateTime)}（{moment(item.CreateTime).fromNow()}）</div>
                    <div style={{marginTop: 5}}>运行耗时： {item.RunTime}</div>
                    <div style={{marginTop: 5}}>返回状态： {item.Status}（非 0 则判定为失败）</div>
                    <div style={{marginTop: 5}}>执行输出： <pre style={preStyle}>{item.Output}</pre></div>
                  </Tabs.TabPane>
                ))}
              </Tabs>
            )}
          </Spin>
        }
        {
          this.props.scheduleInfo == null && 
          <Spin spinning={this.props.scheduleInfoLoading}>
            <Card>
              <p>任务暂无执行记录</p>
            </Card>
          </Spin>
        } 
      </Modal>
    )
  }
}

export default Form.create()(InfoModal)
