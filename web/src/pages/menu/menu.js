import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal,
 Select, Row, Col, Button, Popconfirm, Icon,
 message} from "antd";
import {connect} from "dva";
import {hasPermission} from "@/utils/globalTools";

const FormItem = Form.Item;

@connect(({ loading, menu }) => {
    return {
      menusList: menu.menusList,
      menusLoading: loading.effects['menu/getMenu'],
    }
 })


class MenuPage extends React.Component {
  state = {
    visible: false,
    editCacheData: {},
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({
      type: 'menu/getMenu',
      payload: {
        page: 1,
        pageSize: 10, 
      }
    });
  }
  
  showMenuAddModal = () => {
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
            obj.Icon   === values.Icon 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id;
            dispatch({
              type: 'menu/menuEdit',
              payload: values,
            });
            
          }
        } else {
          values.Type     = 1;
          values.ParentId = 0;
          dispatch({
            type: 'menu/menuAdd',
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
        type: 'menu/menuDel',
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
    values.title =  '编辑菜单-' + values.Name;
    this.setState({ 
      visible: true ,
      editCacheData: values
    });
    
  };

  columns = [
    {
      title: 'ID',
      dataIndex: 'id',
    },
    {
      title: 'Name',
      dataIndex: 'Name',
    },
    {
      title: 'Icon',
      dataIndex: 'Icon',
    },
    {
    title: '操作',
    key: 'action',
    render: (text, record) => (
      <span>
          {hasPermission('menu-edit') && <a onClick={()=>{this.handleEdit(record)}}><Icon type="edit"/>编辑</a>}
        <Divider type="vertical" />
          <Popconfirm title="你确定要删除吗?"  onConfirm={()=>{this.deleteRecord(record.id)}} onCancel={()=>{this.cancel()}}>
            {hasPermission('menu-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
          </Popconfirm>
      </span>
    ),
  },
  ];
  
  render() {
    const {visible, editCacheData} = this.state;
    const {menusList, menusLoading, form: { getFieldDecorator } } = this.props;
    const addmenu = <Button type="primary" onClick={this.showMenuAddModal} >新增一级菜单</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('menu-del') && <Col span={10}>{addmenu}</Col>}
      </Row>;
    return (
      <div>
        <Modal
          title= {editCacheData.title || "新增一级菜单" }
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
          </Form> 
          <FormItem label="图标">
                {getFieldDecorator('Icon', {
                  initialValue: editCacheData.Icon || '',
                  rules: [{ required: true }],
                })(
                  <Input />
                )}
              </FormItem>
        </Modal>

        <Card title="" extra={extra}>
          <Table pagination={false} columns={this.columns} dataSource={menusList} loading={menusLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(MenuPage);