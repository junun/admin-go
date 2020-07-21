import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, 
 Select, Row, Col, Button, Popconfirm, Icon, message, Pagination} from "antd";
import {connect} from "dva";
import {hasPermission} from "@/utils/globalTools";

const FormItem = Form.Item;
const Option = Select.Option;

@connect(({ loading, menu }) => {
    return {
      subMenusLoading: loading.effects['menu/getSubMenu'],
      subMenusList: menu.subMenusList,
      subMenusLen: menu.subMenusLen,
      menusList: menu.menusList,
    }
 })


class SubMenuPage extends React.Component {
  state = {
    visible: false,
    editCacheData: {},
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ type: 'menu/getMenu' });
    dispatch({ 
      type: 'menu/getSubMenu',
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

  // Modal 取消事件
  handleCancel = () => {
    this.setState({
      visible: false,
    });
  };

  // 新增/修改 确定事件
  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        const obj = this.state.editCacheData;
        if (Object.keys(obj).length) {
          if (
            obj.Name      === values.Name && 
            obj.Url       === values.Url && 
            obj.Desc    === values.Desc && 
            obj.Pid === values.Pid && 
            obj.Icon      === values.Icon 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id;
            dispatch({
              type: 'menu/subMenuEdit',
              payload: values,
            });
            
          }
        } else {
          values.type = 1;
          dispatch({
            type: 'menu/subMenuAdd',
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
        type: 'menu/subMenuDel',
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

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'menu/getSubMenu',
      payload: {
        page: page,
        pageSize: 10, 
      }
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
      title: 'Url',
      dataIndex: 'Url',
    },
    {
      title: 'Icon',
      dataIndex: 'Icon',
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
        {hasPermission('submenu-edit') && <a onClick={()=>{this.handleEdit(record)}}><Icon type="edit"/>编辑</a>}
        <Divider type="vertical" />
        <Popconfirm title="你确定要删除吗?"  onConfirm={()=>{this.deleteRecord(record.id)}} onCancel={()=>{this.cancel()}}>
          {hasPermission('submenu-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
        </Popconfirm>
      </span>
    ),
  },
  ];
  
  render() {
    const {visible, editCacheData } = this.state;
    const {menusList, subMenusList, subMenusLen, subMenusLoading, form: { getFieldDecorator } } = this.props;
    const addSubMenu = <Button type="primary" onClick={this.showMenuAddModal} >新增二级菜单</Button>;

    const extra = <Row gutter={16}>
        {hasPermission('submenu-add') && <Col span={10}>{addSubMenu}</Col>}
    </Row>;

    return (
      <div>
        <Modal
          title= { editCacheData.title || "新增二级菜单" }
          destroyOnClose="true"
          // pagination= { {pageSizeOptions: ['30', '40'], showSizeChanger: true}}
          visible= {visible}
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
            <FormItem label="图标">
              {getFieldDecorator('Icon', {
                initialValue: editCacheData.Icon || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="链接">
              {getFieldDecorator('Url', {
                initialValue: editCacheData.Url || '',
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
                {menusList.map(x => <Option key={x.id} value={x.id}>{x.Name}</Option>)}
                </Select>
              )}
            </FormItem>
          </Form> 
        </Modal>

        <Card title="" extra={extra}>
          <Table  
          pagination={{
            showQuickJumper: true,
            total: subMenusLen,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={subMenusList} loading={subMenusLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(SubMenuPage);