import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, 
  Select, Row, Col, Button, Popconfirm, Icon, 
  message} from "antd";
import {connect} from "dva";
import {hasPermission} from "@/utils/globalTools"

const FormItem = Form.Item;
const Option = Select.Option;

@connect(({ loading, menu }) => {
    return {
      subMenusList: menu.subMenusList,
    }
 })

@connect(({ loading, user }) => {
    return {
      permissionsList: user.permissionsList,
      permissionsTotal: user.permissionsTotal,
      permissionsLoading: loading.effects['user/getPermission'],
    }
 })

class PermissionPage extends React.Component {
  state = {
    visible: false,
    editCacheData: {},
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ 
      type: 'menu/getSubMenu',
      payload: {
        page: 1,
        pagesize: 999, 
      }
    });
    dispatch({ 
      type: 'user/getPermission',
      payload: {
        page: 1,
        pagesize: 10, 
      }
    });
  }

  showPremissAddModal = () => {
    this.setState({ 
      visible: true,
      editCacheData: {},
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
            obj.Name       === values.Name && 
            obj.Desc       === values.Desc &&
            obj.Permission === values.Permission &&
            obj.Pid        === values.Pid
          ) {
            message.warning('没有内容修改， 请检查。');
          } else {
            values.id = obj.id;
            dispatch({
              type: 'user/permEdit',
              payload: values,
            });
            
          }
        } else {
          values.type = 2;
          dispatch({
            type: 'user/permAdd',
            payload: values,
          });
        }
        // 重置 `visible` 属性为 false 以关闭对话框
        this.setState({ visible: false });
      }
    });
  };

  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'user/permDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  // Popconfirm 取消事件
  cancel = () => {
  };

  //显示编辑界面
  handleEdit = (values) => {
    values.title =  '编辑权限-' + values.Name;
    this.setState({ 
      visible: true ,
      editCacheData: values
    });
    
  };

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'user/getPermission',
      payload: {
        page: page,
      }
    });
  };

  columns = [
    {
      title: 'Id',
      dataIndex: 'id',
    },
    {
      title: 'Name',
      dataIndex: 'Name',
    },
    {
      title: '权限标识',
      dataIndex: 'Permission',
    },
    {
      title: '描述',
      dataIndex: 'Desc',
    },
    {
      title: '父节点ID',
      dataIndex: 'Pid',
    },
    {
    title: '操作',
    key: 'action',
    render: (text, record) => (
        <span>
          {hasPermission('perm-edit') && <a onClick={()=>{this.handleEdit(record)}}><Icon type="edit"/>编辑</a>}
          <Divider type="vertical" />
          <Popconfirm title="你确定要删除吗?"  onConfirm={()=>{this.deleteRecord(record.id)}} onCancel={()=>{this.cancel()}}>
            {hasPermission('perm-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
          </Popconfirm>
        </span>
      ),
    },
  ];
  
  render() {
    const {visible, editCacheData} = this.state;
    const {subMenusList, permissionsList, permissionsLoading, permissionsTotal, 
      form: { getFieldDecorator } } = this.props;
    const addpremiss = <Button type="primary" onClick={this.showPremissAddModal} >新增权限</Button>;
    const extra = <Row gutter={16}>
        {hasPermission('perm-add') && <Col span={2}>{addpremiss}</Col>}
    </Row>;

    return (
      <div>
        <Modal
          title= { editCacheData.title || "新增权限" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
        >
          <Form>
            <FormItem label="名字">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="权限标识">
              {getFieldDecorator('Permission', {
                initialValue: editCacheData.Permission || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="描述">
              {getFieldDecorator('Desc', {
                initialValue: editCacheData.Desc || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="父菜单">
              {getFieldDecorator('Pid', {
                initialValue: editCacheData.Pid || 'Please select' ,
                rules: [{ required: true }],
              })(
                <Select
                  placeholder="Please select"
                  onChange={this.handleChange}
                  style={{ width: '100%' }}
                >
                  {subMenusList.map(x => <Option key={x.id} value={x.id}>{x.Name}</Option>)}
                </Select>
              )}
            </FormItem>
            
          </Form> 
        </Modal>
        <Card title="" extra={extra}>
          <Table 
          pagination={{
            showQuickJumper: true,
            total: permissionsTotal,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={permissionsList} loading={permissionsLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(PermissionPage);