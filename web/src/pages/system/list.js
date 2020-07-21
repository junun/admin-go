import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, Select,
 Row, Col, Button, Popconfirm, Icon, message} from "antd";
import {connect} from "dva";
import {timeTrans, hasPermission} from "@/utils/globalTools"

const FormItem = Form.Item;
const Option = Select.Option;

@connect(({ loading, user }) => {
    return {
      usersList: user.usersList,
      usersCount: user.usersCount,
      usersLoading: loading.effects['user/getList'],
      rolesList: user.rolesList,
    }
 })

class ListPage extends React.Component {
  state = {
    visible: false,
    editCacheData: {},
    disabled:false,
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ type: 'user/getRole' });
    dispatch({ type: 'user/getList' });
  }
  
  showUserAddModal = () => {
    this.setState({ 
      editCacheData: {},
      visible: true,
    });
  };

  handleCancel = () => {
    this.setState({
      visible: false,
      disabled: false,
    });
  };

  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        const obj = this.state.editCacheData;
        if (Object.keys(obj).length) {
          if (
            obj.Nickname   === values.Nickname && 
            obj.Mobile     === values.Mobile && 
            obj.Email      === values.Email && 
            obj.Rid        === values.Rid 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id;
            values.IsActive =  obj.IsActive;
            dispatch({
              type: 'user/userEdit',
              payload: values,
            });
            
          }
        } else {
          if (values.password != values.repassword) {
            message.warning('两次输入密码不一样！');
            return false;
          } else {
            dispatch({
              type: 'user/userAdd',
              payload: values,
            });
          }
        }
        // 重置 `visible` 属性为 false 以关闭对话框
        this.setState({ 
          visible: false,
        });
      }
    });
  };

  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'user/userDel',
        payload: values.id,
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
    values.title =  '编辑用户-' + values.Nickname;
    this.setState({ 
      visible: true ,
      editCacheData: values,
      disabled:true,
    });
  };

  // 禁用/启用用户
  changeActive = (values) => {
    const { dispatch } = this.props;
    if (values.IsActive) {
      values.IsActive = 0
    } else {
      values.IsActive = 1
    }

    dispatch({
      type: 'user/userEdit',
      payload: values,
    });
  };

  // 重置密码
  restPasswd = (values) => {
    const { dispatch } = this.props;
    values.password = 'ss123456';
    dispatch({
      type: 'user/userEdit',
      payload: values,
    });
  };

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'user/getList',
      payload: {
        page: page,
      }
    });
  };

  columns = [
    {
      title: 'ID',
      dataIndex: 'id',
    },
    {
      title: '姓名',
      dataIndex: 'Nickname',
    },
    {
      title: '登录名',
      dataIndex: 'Name',
    },
    {
      title: '手机',
      dataIndex: 'Mobile',
    },
    {
      title: '邮箱',
      dataIndex: 'Email',
    },
    {
      title: '状态',
      dataIndex: 'IsActive',
      'render': IsActive => 1 && '正常' || '禁用',
    },
    {
      title: '角色',
      dataIndex: 'Rid',
    },
    {
    title: '操作',
    key: 'action',
    render: (text, record) => (
      <span>
        {hasPermission('user-edit') && <a onClick={()=>{this.handleEdit(record)}}><Icon type="edit"/>编辑</a>}
        <Divider type="vertical" />
        <Popconfirm title="你确定要删除吗?"  onConfirm={()=>{this.deleteRecord(record)}} onCancel={()=>{this.cancel()}}>
          {hasPermission('user-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
        </Popconfirm>        
        <Divider type="vertical" /> 
        {
          record.IsActive
          && 
          <Popconfirm title="你确定要禁用用户吗?"  onConfirm={()=>{this.changeActive(record)}} onCancel={()=>{this.cancel()}}>
            {hasPermission('user-edit') && <a title="删除" ><Icon type="lock"/>禁用用户</a>}
          </Popconfirm> 
          || 
          <Popconfirm title="你确定要启用用户吗?"  onConfirm={()=>{this.changeActive(record)}} onCancel={()=>{this.cancel()}}>
            {hasPermission('user-edit') && <a title="删除" ><Icon type="unlock"/>启用用户</a>}
          </Popconfirm> 
        }
        <Divider type="vertical" />  
        <Popconfirm title="你确定要重置吗?"  onConfirm={()=>{this.restPasswd(record)}} onCancel={()=>{this.cancel()}}>
          {hasPermission('user-edit') && <a title="重置" ><Icon type="user"/>重置密码</a>}
        </Popconfirm>
      </span>
    ),
  },
  ];
  
  render() {
    const {visible, editCacheData} = this.state;
    const {rolesList, usersList, usersCount, usersLoading, form: { getFieldDecorator } } = this.props;
    const adduser = <Button type="primary" onClick={this.showUserAddModal} >新增用户</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('user-add') && <Col span={10}>{adduser}</Col>}
      </Row>;


    return (
      <div>
        <Modal
          title= { editCacheData.title || "新建用户" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
        >
          <Form>
            <FormItem label="登录名">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input disabled={this.state.disabled}/>
              )}
            </FormItem>
            <FormItem label="姓名">
              {getFieldDecorator('Nickname', {
                initialValue: editCacheData.Nickname || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="手机">
              {getFieldDecorator('Mobile', {
                initialValue: editCacheData.Mobile || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="邮箱">
              {getFieldDecorator('Email', {
                initialValue: editCacheData.Email || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            { !Object.keys(editCacheData).length  && 
              <FormItem label="密码">
                {getFieldDecorator('password', {
                  initialValue: editCacheData.password || '',
                  rules: [{ required: true }],
                })(
                  <Input type="password" />
                )}
              </FormItem>
            }
            { !Object.keys(editCacheData).length  && 
              <FormItem label="确定密码">
                {getFieldDecorator('repassword', {
                  initialValue: editCacheData.repassword || '',
                  rules: [{ required: true }],
                })(
                  <Input type="password" />
                )}
              </FormItem>
            }
            <FormItem label="角色">
              {getFieldDecorator('Rid', {
                initialValue: editCacheData.Rid || '' ,
                rules: [{ required: true }],
              })(
                <Select
                  placeholder="Please select"
                  onChange={this.handleChange}
                  style={{ width: '100%' }}
                >
                {rolesList.map(x => <Option key={x.id} value={x.id}>{x.Name}</Option>)}
                </Select>
              )}
            </FormItem>
          </Form> 
        </Modal>
        <Card title="" extra={extra}>
          <Table 
          pagination={{
            showQuickJumper: true,
            total: usersCount,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={usersList} loading={usersLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(ListPage);