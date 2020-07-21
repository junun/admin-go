
import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, Alert,
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
    loading: false,
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
    this.setState({loading: true})
    validateFields((err, values) => {
      if (!err) {
        const obj = this.state.editCacheData;
        values.Status  =  values.Status ? 1 : 0
        values.IsCert  =  values.IsCert ? 1 : 0
        if (Object.keys(obj).length) {
          if (
            obj.CertName=== values.CertName &&
            obj.Name    === values.Name     &&
            obj.IsCert  === values.IsCert   &&
            obj.Status  === values.Status   &&
            obj.Desc    === values.Desc 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id;
            dispatch({
              type: 'domain/domainEdit',
              payload: values,
            }).then(()=>{
              this.setState({ visible: false, loading: false})
            })
          }
        } else {
          dispatch({
            type: 'domain/domainAdd',
            payload: values,
          }).then(()=>{
            this.setState({ visible: false, loading: false})
          })
        }
        // 重置 `visible` 属性为 false 以关闭对话框
      }
    });
  };

  //显示编辑界面
  handleEdit = (values) => {
    values.title =  '编辑-' + values.Name;
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
    title: '证书',
    dataIndex: 'CertName',
  }, {
    title: '域名到期日期',
    dataIndex: 'DomainEndTime',
    'render': DomainEndTime => timeDatetimeTrans(DomainEndTime),
  }, {
    title: '证书到期日期',
    dataIndex: 'CertEndTime',
    'render': CertEndTime => timeDatetimeTrans(CertEndTime),
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
    const {visible, editCacheData, loading} = this.state;
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
          confirmLoading={loading}
          width={800}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
        >
          <Alert
            closable
            showIcon
            type="info"
            style={{width: 600, margin: '0 auto 20px', color: '#31708f !important'}}
            message="小提示"
            description={[<div key="1">如果有证书，将会访问你设置的二级域名检测证书有效性信息</div>]} />

          <Form>
            <FormItem label="域名">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input />
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
            <FormItem label="是否有证书">
              {getFieldDecorator('IsCert', {
                initialValue: editCacheData.IsCert && true || false,
                rules: [{ required: true }],
              })(
                <Switch defaultChecked={editCacheData.IsCert && true || false} onChange={this.onCheckChange} />
              )}
            </FormItem>
            <FormItem label="检测证书二级域名">
              {getFieldDecorator('CertName', {
                initialValue: editCacheData.CertName || '',
                rules: [{ required: false }],
              })(
                <Input />
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
