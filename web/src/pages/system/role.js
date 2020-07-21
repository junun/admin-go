import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal,
 Select, Row, Col, Button, Popconfirm, Icon,
 Tree, message, Alert} from "antd";
import {connect} from "dva";
import {timeTrans, hasPermission} from "@/utils/globalTools"
import styles from './role.css';
import PagePerm from "./PagePerm";
import AppPerm from './AppPerm';
import HostPerm from './HostPerm';

const FormItem = Form.Item;
const Option = Select.Option;
const TreeNode = Tree.TreeNode;

@connect(({ loading, user }) => {
    return {
      rolesList: user.rolesList,
      rolesCount: user.rolesCount,
      usersLoading: loading.effects['user/getRole'],
      allPermissionsList: user.allPermissionsList,
      userPermissionsList: user.userPermissionsList,
      roleVisible: user.roleVisible,
      // checkedKeys: [],
    }
 })

class RolePage extends React.Component {
  state = {
    roleVisible: this.props.roleVisible,
    editCacheData: {},
    value: '用户列表',
    roleId: '',
    checkedKeys: [],
    pagePermVisible: false,
    appPermVisible: false,
    hostPermVisible: false,
    rid:0,
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ type: 'user/getRole' });
    dispatch({ type: 'user/getAllPermission' });
  }
  
  showRoleAddModal = () => {
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
            obj.Name       === values.Name && 
            obj.Desc     === values.Desc 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id;
            dispatch({
              type: 'user/roleEdit',
              payload: values,
            });
            
          }
        } else {
          dispatch({
            type: 'user/roleAdd',
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
        type: 'user/roleDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  // Popconfirm 取消事件
  cancel = () => {
  };

  switchPagePerm = () => {
    this.setState({ 
      pagePermVisible: !this.state.pagePermVisible,
    });
  };

  switchAppPerm = () => {
    this.setState({ 
      appPermVisible: !this.state.appPermVisible,
    });
  }

  switchHostPerm = () => {
    this.setState({ 
      hostPermVisible: !this.state.hostPermVisible,
    });
  }

  //显示编辑界面
  handleEdit = (values) => {
    values.title =  '编辑角色-' + values.Name;
    this.setState({ 
      visible: true ,
      editCacheData: values
    });
    
  };

  //显示权限界面
  handlePermission = (values) => {
    values.title =  '功能权限-' + values.Name;
    this.setState({ 
      editCacheData: values,
      pagePermVisible: true,
      rid: values.id,
    });
  };

  //取消权限界面
  handlePermissionCancel = () => {
    const { dispatch } = this.props;
    dispatch({ 
      type: 'user/cancelUserPermission',
    });
  };

  handleAppPerm = (values) => {
    values.title =  '应用权限-' + values.Name;
    this.setState({ 
      editCacheData: values,
      appPermVisible: true,
      rid: values.id,
    });
  }

  handleHostPerm = (values) => {
    values.title =  '主机权限-' + values.Name;
    this.setState({ 
      editCacheData: values,
      hostPermVisible: true,
      rid: values.id,
    });
  }

  onCheck = (values) => {
    this.setState({ 
      checkedKeys: values,
    });
  };

  handlePermissionOk = (item) => {
    console.log(item);
    const { dispatch, userPermissionsList } = this.props;
    // const keys = this.state.checkedKeys;
    const keys = item;
    if (keys.length >0) {
      const values = {};
      values.id = this.state.editCacheData.id;
      values.Codes = keys;
      dispatch({
        type: 'user/rolePermsAdd',
        payload: values,
      });
    } else {
      message.warning('没有内容修改， 请检查。');
      return false;
    }
    
    dispatch({ 
      type: 'user/cancelUserPermission',
    });
  };

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'user/getRole',
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
      title: '角色名',
      dataIndex: 'Name',
    },
    {
      title: '描述',
      dataIndex: 'Desc',
    },
    {
      title: '操作',
      key: 'action',
      render: (text, record) => (
        <span>
          {hasPermission('role-edit') && <a onClick={()=>{this.handleEdit(record)}}><Icon type="edit"/>编辑</a>}
          <Divider type="vertical" />
          {
            hasPermission('role-perm-list') && 
            <a onClick={()=>{this.handlePermission(record)}}>
              <Icon type="tags"/>功能权限
            </a>
          }
          <Divider type="vertical" />
          {
            hasPermission('role-perm-list') && 
            <a onClick={()=>{this.handleAppPerm(record)}}>
              <Icon type="ci"/>发布权限
            </a>
          }
          <Divider type="vertical" />
          {
            hasPermission('role-perm-list') && 
            <a onClick={()=>{this.handleHostPerm(record)}}>
              <Icon type="desktop"/>主机权限
            </a>
          }
          <Divider type="vertical" />
          <Popconfirm title="你确定要删除吗?"  onConfirm={()=>{this.deleteRecord(record)}} onCancel={()=>{this.cancel()}}>
            {hasPermission('role-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
          </Popconfirm>
        </span>
      ),
    },
  ];
  
  render() {
    const {pagePermVisible, appPermVisible, hostPermVisible,
      value, visible, editCacheData} = this.state;
    const {userPermissionsList, roleVisible,
     allPermissionsList, rolesList, rolesCount,
    rolesLoading, form: { getFieldDecorator } } = this.props;

    const addrole = <Button type="primary" onClick={this.showRoleAddModal} >新增角色</Button>;

    const extra = <Row gutter={16}>
        {hasPermission('role-add') && <Col span={10}>{addrole}</Col>}
    </Row>;
    
    const loop = data => data.map((item) => {
      if (item.children && item.children.length) {
        return <TreeNode key={item.id} title={item.Name} value={item.id} >{loop(item.children)}</TreeNode>;
      }
      return <TreeNode value={item.id} key={item.id} title={item.Name} />;
    });

    return (
      <div>
        <Modal
          title= { editCacheData.title  }
          visible= {roleVisible}
          destroyOnClose= "true"
          // onOk={this.handlePermissionOk}
          onCancel={this.handlePermissionCancel}
        >
          <Tree 
            showLine
            checkable
            multiple={true}
            defaultExpandAll={true}
            defaultCheckedKeys={userPermissionsList || ''}
            onSelect={this.onSelect} 
            onCheck={this.onCheck}
          >
            {loop(allPermissionsList)}
          </Tree>
        </Modal>
        
        <Modal
          title= { editCacheData.title || "新增角色" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
        >
          <Form>
            <FormItem label="角色名">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
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
          </Form> 
        </Modal>

        {pagePermVisible && 
          <PagePerm 
            rid={this.state.rid} 
            allPerm={allPermissionsList} 
            rolePerm={userPermissionsList} 
            dispatch={this.props.dispatch}
            onCancel={this.switchPagePerm} 
            onOk={this.handlePermissionOk}
          />
        }

        {appPermVisible &&
          <AppPerm
            rid={this.state.rid} 
            dispatch={this.props.dispatch}
            onCancel={this.switchAppPerm} 
          />
        }

        {hostPermVisible &&
          <HostPerm
            rid={this.state.rid} 
            dispatch={this.props.dispatch}
            onCancel={this.switchHostPerm} 
          />
        }

        <Card title="" extra={extra}>
          <Table 
          pagination={{
            showQuickJumper: true,
            total: rolesCount,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={rolesList} loading={rolesLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(RolePage);