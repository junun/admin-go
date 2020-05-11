
import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal,
 Select, Row, Col, Button, Popconfirm, Icon, 
 message} from "antd";
 import {hasPermission} from "@/utils/globalTools"

import {connect} from "dva";

const FormItem = Form.Item;

@connect(({ loading, config }) => {
    return {
      appTypeList: config.appTypeList,
      appTypeLen: config.appTypeLen,
      appTypeLoading: loading.effects['config/getAppType'],
    }
})

class AppTypePage extends React.Component {
  state = {
    visible: false,
    editCacheData: {},
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({
      type: 'config/getAppType',
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
      type: 'config/getAppType',
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
        if (Object.keys(obj).length) {
          if (
            obj.Name   === values.Name && 
            obj.Desc   === values.Desc 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id;
            dispatch({
              type: 'config/appTypeEdit',
              payload: values,
            });
            
          }
        } else {
          dispatch({
            type: 'config/appTypeAdd',
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
    values.title =  '编辑菜单-' + values.name;
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
        type: 'config/appTypeDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  columns = [
  {
    title: '应用类别',
    dataIndex: 'Name',
  }, {
    title: '备注',
    dataIndex: 'Desc',
    ellipsis: true
  }, {
    title: '操作',
    width: 200,
    render: (text, record) => (
      <span>
        {hasPermission('config-env-edit') && <a onClick={()=>{this.handleEdit(record)}}><Icon type="edit"/>编辑</a>}
        <Divider type="vertical" />
        {
          hasPermission('config-env-del') && 
          <Popconfirm title="你确定要删除吗?"  onConfirm={()=>{this.deleteRecord(record.id)}} onCancel={()=>{this.cancel()}}>
            <a title="删除" ><Icon type="delete"/>删除</a>
          </Popconfirm>
        }
      </span>
    ),
  }];
  
  render() {
    const {visible, editCacheData} = this.state;
    const {appTypeList, appTypeLen, appTypeLoading, form: { getFieldDecorator } } = this.props;
    const addappType = <Button type="primary" onClick={this.showTypeAddModal} >新增应用类别</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('config-env-add') && <Col span={10}>{addappType}</Col>}
      </Row>;
    return (
      <div>
        <Modal
          title= {editCacheData.title || "新增应用类别" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
        >
          <Form>
            <FormItem label="类型名字">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            
            <FormItem label="备注信息">
              {getFieldDecorator('Desc', {
                initialValue: editCacheData.Desc || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
          </Form> 
        </Modal>

        <Card title="" extra={extra}>
          <Table  
          pagination={{
            showQuickJumper: true,
            total: appTypeLen,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={appTypeList} loading={appTypeLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(AppTypePage);
