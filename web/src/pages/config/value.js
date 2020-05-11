
import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal,
 Select, Row, Col, Button, Popconfirm, Icon, 
 message} from "antd";
import {hasPermission} from "@/utils/globalTools"
import {connect} from "dva";

const FormItem = Form.Item;

@connect(({ loading, config}) => {
    return {
      appValueList: config.appValueList,
      appValueLen: config.appValueLen,
      appValueLoading: loading.effects['config/getAppValue'],
      projectsList: config.projectsList,
    }
})

class AppValuePage extends React.Component {
  state = {
    editCacheData: [],
    visible: false,
    aid : 0,
  };

  componentDidMount() {
    const { location, dispatch } = this.props;
    this.setState({ 
      aid: location.search.split("=")[1] 
    });

    dispatch({
      type: 'config/getAppValue',
      payload: {
        page: 1,
        pageSize: 10,
        aid: location.search.split("=")[1],
      }
    });
  }

  // 翻页
  pageChange = (page) => {
    const { dispatch, location } = this.props;
    dispatch({
      type: 'config/getAppValue',
      payload: {
        page: page,
        pageSize: 10,
        aid: location.search.split("=")[1],
      }
    });
  };

  showTypeAddModal = () => {
    this.setState({ 
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
            obj.Value  === values.Value && 
            obj.Desc   === values.Desc 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id  = obj.id;
            values.Aid = parseInt(this.state.aid)
            dispatch({
              type: 'config/appValueEdit',
              payload: values,
            });
            
          }
        } else {
          values.Aid = parseInt(this.state.aid)
          dispatch({
            type: 'config/appValueAdd',
            payload: values,
          });
        }
        // 重置 `visible` 属性为 false 以关闭对话框
        this.setState({ visible: false });
      }
    });
  };

  cancel  = () => { 
  }

  //显示编辑界面
  handleEdit = (values, page) => {
    values.title =  '编辑应用变量-' + values.Name;
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
        type: 'config/appValueDel',
        payload: values,
      });
    } else {
      message.error('错误的id');
    }
  };

  columns = [
    {
      title: 'ID',
      dataIndex: 'id',
    },
    {
      title: '变量名',
      dataIndex: 'Name',
    },
    {
      title: '变量值',
      dataIndex: 'Value',
    },
    {
      title: '备注',
      dataIndex: 'Desc',
      ellipsis: true
    }, {
      title: '操作',
      width: 200,
      render: (text, record) => (
        <span>
          {
            hasPermission('app-value-edit') && 
            <a onClick={()=>{this.handleEdit(record)}}>
              <Icon type="edit"/>编辑
            </a>
          }
          <Divider type="vertical" />
          <Popconfirm title="你确定要删除吗?" 
            onConfirm={()=>{this.deleteRecord(record.id)}} 
            onCancel={this.cancel()}
          >
            {hasPermission('app-value-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
          </Popconfirm>
        </span>
      ),
  }];
  
  render() {
    const {visible, aid, editCacheData} = this.state;
    const {projectsList, appValueList, appValueLen, appValueLoading, form: { getFieldDecorator } } = this.props;
    const addAppValue = <Button type="primary" onClick={this.showTypeAddModal} >新增变量</Button>;

    const extra = <Row gutter={16}>
          {hasPermission('app-value-add') && <Col span={10}>{addAppValue}</Col>}
      </Row>;
    return (
      <div>
        <Modal
          title= { editCacheData.title || "新增变量" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
        >
          <Form>
            <FormItem label="变量名">
              {getFieldDecorator('Name', {
                initialValue: editCacheData.Name || '',
                rules: [{ required: true }],
              })(
                <Input />
              )}
            </FormItem>
            <FormItem label="变量值">
              {getFieldDecorator('Value', {
                initialValue: editCacheData.Value || '',
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
                <Input.TextArea />
              )}
            </FormItem>
          </Form> 
        </Modal>

        <Card title="" extra={extra}>
          <Table  
          pagination={{
            showQuickJumper: true,
            total: appValueLen,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={appValueList} loading={appValueLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(AppValuePage);