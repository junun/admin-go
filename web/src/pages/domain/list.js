
import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal,
 Select, Row, Col, Button, Popconfirm, Icon, Switch, DatePicker,
 message} from "antd";
import {hasPermission, timeDatetimeTrans} from "@/utils/globalTools"
import {connect} from "dva";
import moment from 'moment';

const dateFormat = 'YYYY-MM-DD';
const FormItem = Form.Item;

@connect(({ loading, domain }) => {
    return {
      domainList: domain.domainList,
      domainListLen: domain.domainListLen,
      domainLoading: loading.effects['domain/getDomain'],
    }
})


class DomainPage extends React.Component {
  state = {
    visible: false,
    editCacheData: {},
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({
      type: 'domain/getDomain',
      payload: {
        page: 1,
        pageSize: 10, 
      }
    });
  }

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'domain/getDomain',
      payload: {
        page: page,
        pageSize: 10, 
      }
    });
  };

  showTypeAddModal = () => {
    this.setState({ 
      editCacheData: {},
      visible: true 
    });
  };

  handleCancel = () => {
    this.setState({
      visible: false,
    });
  };

  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        const obj = this.state.editCacheData;
        values.Status  =  values.Status ? 1 : 0
        if (Object.keys(obj).length) {
          if (
            obj.Name    === values.Name &&
            obj.Channel === values.Channel &&
            obj.EndTime === values.EndTime &&
            obj.Status  === values.Status &&
            obj.Desc    === values.Desc 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id;
            dispatch({
              type: 'domain/domainEdit',
              payload: values,
            });
            
          }
        } else {
          values.Status = values.Status ? 1 : 0
          dispatch({
            type: 'domain/domainAdd',
            payload: values,
          });
        }
        // 重置 `visible` 属性为 false 以关闭对话框
        this.setState({ visible: false });
      }
    });
  };

  //显示编辑界面
  handleEdit = (values) => {
    values.title =  '编辑-' + values.name;
    this.setState({ 
      visible: true ,
      editCacheData: values
    });
  };

  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'domain/domainDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  columns = [
  {
    title: '域名',
    dataIndex: 'Name',
  }, {
    title: '申请渠道',
    dataIndex: 'Channel',
    ellipsis: true
  }, {
    title: '申请日期',
    dataIndex: 'StartTime',
    'render': StartTime => timeDatetimeTrans(StartTime),
  }, {
    title: '到期日期',
    dataIndex: 'EndTime',
    'render': EndTime => timeDatetimeTrans(EndTime),
  }, {
    title: '状态',
    dataIndex: 'Status',
    'render': Status => Status && '使用中' || '已经过期',
  }, {
    title: '备注',
    dataIndex: 'Desc',
    ellipsis: true
  }, {
    title: '操作',
    width: 200,
    render: (text, record) => (
      <span>
          {hasPermission('domain-info-edit') && <a onClick={()=>{this.handleEdit(record)}}><Icon type="edit"/>编辑</a>}
        <Divider type="vertical" />
          <Popconfirm title="你确定要删除吗?"  onConfirm={()=>{this.deleteRecord(record.id)}} onCancel={()=>{this.cancel()}}>
            {hasPermission('domain-info-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
          </Popconfirm>
      </span>
    ),
  }];
  
  render() {
    const {visible, editCacheData} = this.state;
    const {domainList, domainListLen, domainLoading, form: { getFieldDecorator } } = this.props;
    const adddomainRple = <Button type="primary" onClick={this.showTypeAddModal} >新增域名</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('domain-info-add') && <Col span={10}>{adddomainRple}</Col>}
      </Row>;
    return (
      <div>
        <Modal
          title= {editCacheData.title || "新增域名" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}>
          <Form>
            <FormItem label="域名">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="申请渠道">
              {getFieldDecorator('Channel', {
                initialValue: editCacheData.Channel || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>

            <FormItem label="申请日期">
              {getFieldDecorator('StartTime', {
                initialValue: editCacheData.StartTime && moment(editCacheData.StartTime, dateFormat) || null,
                rules: [{ required: true }],
              })(
                <DatePicker onChange={this.onCheckChange} />
              )}
            </FormItem>

            <FormItem label="到期日期">
              {getFieldDecorator('EndTime', {
                initialValue: editCacheData.EndTime && moment(editCacheData.EndTime, dateFormat) || null,
                rules: [{ required: true }],
              })(
                <DatePicker onChange={this.onCheckChange} />
              )}
            </FormItem>
            <FormItem label="是否有效">
              {getFieldDecorator('Status', {
                initialValue: editCacheData.Status && true || false,
                rules: [{ required: true }],
              })(
                <Switch defaultChecked={editCacheData.Status && true || false} onChange={this.onCheckChange} />
              )}
            </FormItem>
            <FormItem label="备注信息">
              {getFieldDecorator('Desc', {
                initialValue: editCacheData.Desc || '',
                rules: [{ required: true }],
              })(
                <Input.TextArea />
              )}
            </FormItem>
          </Form> 
        </Modal>

        <Card title="" extra={extra}>
          <Table  
          pagination={{
            showQuickJumper: true,
            total: domainListLen,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={domainList} loading={domainLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(DomainPage);
