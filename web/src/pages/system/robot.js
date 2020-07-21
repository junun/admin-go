import React, {Fragment, Component} from "react";
import {Form, Card, Input, Table, Divider, Modal, Select,
 Row, Col, Button, Popconfirm, Icon, message} from "antd";
import {connect} from "dva";
import {timeTrans, hasPermission} from "@/utils/globalTools"

const FormItem = Form.Item;
const Option = Select.Option;

const robotChannel = [
  {id: 1, name: '钉钉数字签名机器人'},
  {id: 2, name: '钉钉关键字机器人'},
  {id: 3, name: '钉钉Acl机器人'},
  {id: 4, name: '企业微信应用'},
  {id: 5, name: '企业微信机器人'},
];

@connect(({ loading, user }) => {
    return {
      robotList: user.robotList,
      robotCount: user.robotCount,
      robotLoading: loading.effects['user/getRobot'],
    }
 })

class RobotPage extends React.Component {
  state = {
    visible: false,
    editCacheData: {},
    showType: 0,
  };

  componentDidMount() {
    const { dispatch } = this.props;
    dispatch({ 
      type: 'user/getRobot',
      payload: {
        page: 1,
        pageSize: 10, 
      }
    });
  }
  
  showAddModal = () => {
    this.setState({ 
      editCacheData: {},
      visible: true,
      showType: 0,
    });
  };

  handleCancel = () => {
    this.setState({
      visible: false,
      showType: 0,
    });
  };

  handleOk = () => {
    const { dispatch, form: { validateFields } } = this.props;
    validateFields((err, values) => {
      if (!err) {
        const obj = this.state.editCacheData;
        if (Object.keys(obj).length) {
          if (
            obj.Name     === values.Name && 
            obj.Webhook  === values.Webhook && 
            obj.Select   === values.Select && 
            obj.Keyword  === values.Keyword && 
            obj.Desc     === values.Desc && 
            obj.Type     === values.Type 
          ) {
            message.warning('没有内容修改， 请检查。');
            return false;
          } else {
            values.id = obj.id
            values.Status = obj.Status
            dispatch({
              type: 'user/robotEdit',
              payload: values,
            });
          }
        } else {
          values.Status = 1;
          dispatch({
            type: 'user/robotAdd',
            payload: values,
          });
        }
        // 重置 `visible` 属性为 false 以关闭对话框
        this.setState({ 
          visible: false,
          showType: 0,
        });
      }
    });
  };

  // 删除一条记录
  deleteRecord = (values) => {
    const { dispatch } = this.props;
    if (values) {
      dispatch({
        type: 'user/robotDel',
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
    values.title =  '编辑-' + values.Name;
    this.setState({ 
      showType: values.Type,
      visible: true ,
      editCacheData: values,
    });
  }

  handleTest = (values) => {
    this.setState({loading: true});
    this.props.dispatch({
      type: 'user/robotTest',
      payload: values,
    }).finally(() => this.setState({loading: false}))
  }

  // 禁用/启用
  changeActive = (values) => {
    const { dispatch } = this.props;
    if (values.Status) {
      values.Status = 0
    } else {
      values.Status = 1
    }

    dispatch({
      type: 'user/robotEdit',
      payload: values,
    });
  };

  // 翻页
  pageChange = (page) => {
    const { dispatch } = this.props;
    dispatch({
      type: 'user/getRobot',
      payload: {
        page: page,
      }
    });
  }

  handleChange = (e) => {
    this.setState({showType: e});
  }

  columns = [
    {
      title: 'ID',
      dataIndex: 'id',
    },
    {
      title: '标识',
      dataIndex: 'Name',
    },
    {
      title: '通道',
      dataIndex: 'Type',
      'render': Type => robotChannel.map(x => {
        if (Type == x.id) {
          return x.name
        }
      })
    },
    {
      title: '状态',
      dataIndex: 'Status',
      'render': Status => Status == 1 && '正常' || '禁用',
    },
    {
      title: '备注',
      dataIndex: 'Desc',
    },
    {
      title: '操作',
      key: 'action',
      render: (record) => (
        <span>
          {hasPermission('setting-robot-edit') && <a onClick={()=>{this.handleEdit(record)}}><Icon type="edit"/>编辑</a>}
          <Divider type="vertical" />
          <Popconfirm title="你确定要删除吗?"  onConfirm={()=>{this.deleteRecord(record)}} onCancel={()=>{this.cancel()}}>
            {hasPermission('setting-robot-del') && <a title="删除" ><Icon type="delete"/>删除</a>}
          </Popconfirm>        
          <Divider type="vertical" /> 
          {
            record.Status
            && 
            <Popconfirm title="你确定要禁用该通道吗?"  onConfirm={()=>{this.changeActive(record)}} onCancel={()=>{this.cancel()}}>
              {hasPermission('setting-robot-edit') && <a><Icon type="lock"/>禁用</a>}
            </Popconfirm> 
            || 
            <Popconfirm title="你确定要启用该通道吗?"  onConfirm={()=>{this.changeActive(record)}} onCancel={()=>{this.cancel()}}>
              {hasPermission('setting-robot-edit') && <a><Icon type="unlock"/>启用</a>}
            </Popconfirm> 
          }
          <Divider type="vertical" /> 
          {hasPermission('setting-robot-test') && <a disabled={!record.Status} loading={this.state.loading} onClick={()=>{this.handleTest(record)}}><Icon type="check"/>测试</a>}
        </span>
      ),
    },
  ];
  
  render() {
    const {visible, editCacheData, showType} = this.state;
    const {robotList, robotCount, robotLoading, form: { getFieldDecorator } } = this.props;
    const AddRobot = <Button type="primary" onClick={this.showAddModal} >新增机器人</Button>;
    const extra = <Row gutter={16}>
          {hasPermission('setting-robot-add') && <Col span={10}>{AddRobot}</Col>}
      </Row>;

    return (
      <div>
        <Modal
          title= { editCacheData.title || "新增机器人" }
          visible= {visible}
          destroyOnClose= "true"
          onOk={this.handleOk}
          onCancel={this.handleCancel}
        >
          <Form>
            <FormItem label="通道类型">
              {getFieldDecorator('Type', {
                initialValue: editCacheData.Type || '' ,
                rules: [{ required: true }],
              })(
                <Select
                  placeholder="Please select"
                  onChange={e => this.handleChange(e)}
                  style={{ width: '100%' }}
                >
                {robotChannel.map(x => <Select.Option key={x.id} value={x.id}>{x.name}</Select.Option>)}
                </Select>
              )}
            </FormItem>
            { showType != 0 &&
              <FormItem label="唯一标识">
                {getFieldDecorator('Name', {
                  initialValue: editCacheData.Name || '',
                  rules: [{ required: true }],
                })(
                  <Input/>
                )}
              </FormItem>
            }
            { showType != 0 &&  showType != 4 &&
              <FormItem label="Webhook">
                {getFieldDecorator('Webhook', {
                  initialValue: editCacheData.Webhook || '',
                  rules: [{ required: true }],
                })(
                  <Input.TextArea />
                )}
              </FormItem>
            }
            { showType == 1 &&  
              <FormItem label="Secret">
                {getFieldDecorator('Secret', {
                  initialValue: editCacheData.Secret || '',
                  rules: [{ required: true }],
                })(
                  <Input />
                )}
              </FormItem>
            }
            { showType == 2 && 
              <FormItem label="关键字">
                {getFieldDecorator('Keyword', {
                  initialValue: editCacheData.Keyword || '',
                  rules: [{ required: true }],
                })(
                  <Input />
                )}
              </FormItem>
            }
            { showType == 4 &&
              <FormItem label="corpid" help="企业拥有唯一的corpid https://work.weixin.qq.com/api/doc/90000/90135/90665#corpid">
                {getFieldDecorator('Webhook', {
                  initialValue: editCacheData.Webhook || '',
                  rules: [{ required: true }],
                })(
                  <Input />
                )}
              </FormItem>
            }
            { showType == 4 &&
              <FormItem label="agentid" help="应用id 参考https://work.weixin.qq.com/api/doc/90000/90135/90665#corpid">
                {getFieldDecorator('Keyword', {
                  initialValue: editCacheData.Keyword || '',
                  rules: [{ required: true }],
                })(
                  <Input />
                )}
              </FormItem>
            }
            { showType == 4 &&
              <FormItem label="secret" help="应用secret 参考https://work.weixin.qq.com/api/doc/90000/90135/90665#corpid">
                {getFieldDecorator('Secret', {
                  initialValue: editCacheData.Secret || '',
                  rules: [{ required: true }],
                })(
                  <Input />
                )}
              </FormItem>
            }
            { showType != 0 &&
              <FormItem label="备注信息">
                {getFieldDecorator('Desc', {
                  initialValue: editCacheData.Desc || '',
                  rules: [{ required: false }],
                })(
                  <Input.TextArea />
                )}
              </FormItem>
            }
          </Form> 
        </Modal>

        <Card title="" extra={extra}>
          <Table 
          pagination={{
            showQuickJumper: true,
            total: robotCount,
            showTotal: (total, range) => `第${range[0]}-${range[1]}条 总共${total}条`,
            onChange: this.pageChange
          }}
          columns={this.columns} dataSource={robotList} loading={robotLoading} rowKey="id" />
        </Card>
      </div>
    );
  }
}

export default Form.create()(RobotPage);